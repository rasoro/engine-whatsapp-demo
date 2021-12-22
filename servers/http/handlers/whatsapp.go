package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/services"
)

const confirmationMessage = "Token válido, Whatsapp demo está pronto para sua utilização"

type WhatsappHandler struct {
	ContactService  services.ContactService
	ChannelService  services.ChannelService
	CourierService  services.CourierService
	WhatsappService services.WhatsappService
}

func (h *WhatsappHandler) HandleIncomingRequests(w http.ResponseWriter, r *http.Request) {
	incomingMsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("unexpected server error - " + err.Error())
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	incomingContact := parseToContact(string(incomingMsg))
	if incomingContact == nil {
		err := errors.New("request without being from a contact")
		logger.Debug(fmt.Sprintf("%v: %v", err.Error(), string(incomingMsg)))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	contact, err := h.ContactService.FindContact(incomingContact)
	if err != nil {
		logger.Debug(err.Error())
	}

	if possibleToken := extractTextMessage(string(incomingMsg)); possibleToken != "" {
		channelFromToken, err := h.ChannelService.FindChannelByToken(possibleToken)
		if err != nil {
			logger.Error(err.Error())
		}
		if channelFromToken != nil {
			incomingContact.Channel = channelFromToken.ID
			if contact != nil {
				contact.Channel = channelFromToken.ID
				_, err = h.ContactService.UpdateContact(contact)
				if err != nil {
					logger.Error(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, b, err := h.sendTokenConfirmation(contact)
				if err != nil {
					logger.Error(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					body, _ := ioutil.ReadAll(b)
					logger.Debug(string(body))
					w.WriteHeader(http.StatusOK)
					return
				}
			} else {
				_, err := h.ContactService.CreateContact(incomingContact)
				if err != nil {
					logger.Error(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, b, err := h.sendTokenConfirmation(incomingContact)
				if err != nil {
					logger.Error(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					body, _ := ioutil.ReadAll(b)
					logger.Debug(string(body))
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		} else {
			if contact != nil {
				channelId := contact.Channel.Hex()
				channel, err := h.ChannelService.FindChannelById(channelId)
				if err != nil {
					logger.Debug(err.Error())
				}
				if channel != nil {
					channelUUID := channel.UUID
					status, err := h.CourierService.RedirectMessage(channelUUID, string(incomingMsg))
					if err != nil {
						logger.Debug(err.Error())
						w.WriteHeader(status)
						fmt.Fprint(w, err)
					}
				}
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, errors.New("contact not found and token not valid"))
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	wconfig := config.GetConfig().Whatsapp
	httpClient := &http.Client{}
	reqPath := "/v1/users/login"

	reqURL, _ := url.Parse(wconfig.BaseURL + reqPath)

	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{},
		Body:   r.Body,
	}

	req.SetBasicAuth(config.AppConf.Whatsapp.Username, config.AppConf.Whatsapp.Password)

	res, err := httpClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}

	var login struct {
		Users []struct {
			Token        string
			ExpiresAfter string
		}
		Meta struct {
			Version   string
			ApiStatus string
		}
	}

	bdBytes, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	bdString := string(bdBytes)
	log.Println(bdString)

	if err := json.Unmarshal(bdBytes, &login); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}

	newToken := login.Users[0].Token

	config.UpdateToken(newToken)
	logger.Info("Whatsapp token update")
	w.WriteHeader(http.StatusOK)
	for k, v := range res.Header {
		w.Header().Set(k, strings.Join(v, ""))
	}
	fmt.Fprint(w, bdString)
}

func (h *WhatsappHandler) sendTokenConfirmation(contact *models.Contact) (http.Header, io.ReadCloser, error) {
	urn := contact.URN
	payload := fmt.Sprintf(
		`{"to":"%s","type":"text","text":{"body":"%s"}}`,
		urn,
		confirmationMessage,
	)
	payloadBytes := []byte(payload)

	return h.WhatsappService.SendMessage(payloadBytes)
}

func parseToContact(m string) *models.Contact {
	name := extractName(m)
	number := extractNumber(m)
	if name != "" && number != "" {
		return &models.Contact{
			URN:  number,
			Name: name,
		}
	}
	return nil
}

func extractName(m string) string {
	var result map[string][]map[string]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["contacts"] != nil {
		return result["contacts"][0]["profile"]["name"].(string)
	}
	return ""
}

func extractNumber(m string) string {
	var result map[string][]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["messages"] != nil {
		return result["messages"][0]["from"].(string)
	}
	return ""
}

func extractTextMessage(m string) string {
	var result map[string][]map[string]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["messages"] != nil {
		return result["messages"][0]["text"]["body"].(string)
	}
	return ""
}
