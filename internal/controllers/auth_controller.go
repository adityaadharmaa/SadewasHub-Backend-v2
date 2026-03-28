package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sadewashub-go/internal/config"
	"sadewashub-go/internal/models"
	"sadewashub-go/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
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
	if c.Query("state") != "random-state" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State tidak valid"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code tidak ditemukan"})
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menukar token"})
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca stream body"})
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "Gagal membaca JSON data user",
			"pesan_asli_google": string(bodyBytes),
			"error_detail":      err.Error(),
		})
		return
	}

	if _, ok := userInfo["error"]; ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Google menolak akses",
			"detail": userInfo,
		})
		return
	}

	email, _ := userInfo["email"].(string)
	name, _ := userInfo["name"].(string)
	googleID, _ := userInfo["id"].(string)
	picture, _ := userInfo["picture"].(string)

	var user models.User

	result := config.DB.Where("email = ? ", email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			userID := uuid.New()
			now := time.Now()
			socialType := "google"

			user = models.User{
				ID:              userID,
				Email:           email,
				EmailVerifiedAt: &now,
				SocialID:        &googleID,
				SocialType:      &socialType,
			}

			profile := models.UserProfile{
				ID:       uuid.New(),
				UserID:   userID,
				FullName: name,
			}

			profilePic := models.Attachment{
				ID:             uuid.New(),
				AttachableID:   userID,
				AttachableType: "App\\Models\\User",
				FilePath:       picture,
				FileType:       "avatar_google",
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			var tenantRole models.Role
			if err := config.DB.Where("name = ?", "tenant").First(&tenantRole).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Role tenant belum dibuat di database"})
				return
			}

			tx := config.DB.Begin()
			if err := tx.Create(&user).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan user baru"})
				return
			}
			if err := tx.Create(&profile).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profile user"})
				return
			}
			if err := tx.Create(&profilePic).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan foto profile"})
				return
			}
			if err := tx.Model(&user).Association("Roles").Append(&tenantRole); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memberikan role tenant ke user"})
				return
			}
			tx.Commit()
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		if user.SocialID == nil {
			socialType := "google"
			config.DB.Model(&user).Updates(models.User{SocialID: &googleID, SocialType: &socialType})
		}
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat jwt token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil",
		"token":   jwtToken,
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"name":    name,
			"picture": picture,
		},
	})
}
