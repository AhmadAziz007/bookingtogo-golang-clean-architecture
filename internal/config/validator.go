package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("valid_date", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()

		// ✅ Jika kosong, anggap valid — biarkan tag "required" yang handle wajib/tidaknya
		if dateStr == "" {
			return true
		}

		_, err := time.Parse("2006-01-02", dateStr)
		return err == nil
	})

	return validate
}
