package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func InitOAuth() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// Endpoint untuk memulai google login
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("random-state")

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Endpoint Callback dari google
func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "random-state" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State tidak valid"})
		return
	}

	// Ambil "code" yang dikirim google
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code tidak ditemukan"})
		return
	}

	// Tukar "code" dengan Access Token
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menukar token: " + err.Error()})
		return
	}

	// Gunakan akses token untuk mengambil data profile user dari Google API
	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user: " + err.Error()})
		return
	}

	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"message":          "Berhasil mendapatkan data dari google",
		"google_user_info": userInfo,
	})
}
