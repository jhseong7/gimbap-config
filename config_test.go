package config_test

import (
	"os"
	"path/filepath"

	config "github.com/jhseong7/gimbap-config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("ConfigService", func() {
	Context("When the configuration file exists", func() {
		AfterEach(func() {
			viper.Reset()
		})

		It("Check string and numeric value read", func() {
			svc := config.NewConfigService(config.ConfigOption{
				ConfigFilePathList: []string{"test/.env.test", "test/test.yaml"},
			})

			Expect(svc.GetString("env_string")).To(Equal("test"))
			Expect(svc.GetString("ENV_STRING")).To(Equal("test"))
			Expect(svc.GetInt("env_number")).To(Equal(1))
			Expect(svc.GetString("yaml_str")).To(Equal("string"))
			Expect(svc.GetInt("yaml_int")).To(Equal(3120))
			Expect(svc.GetFloat64("yaml_float")).To(Equal(3.14))
		})

		It("Read from raw configuration data", func() {
			svc := config.NewConfigService(config.ConfigOption{
				ConfigData: map[string]interface{}{
					"raw_string": "test",
					"raw_number": 1,
				},
			})

			Expect(svc.GetString("raw_string")).To(Equal("test"))
			Expect(svc.GetInt("raw_number")).To(Equal(1))
		})

		It("Check the configuration file change", func() {
			// Create a temporary configuration file
			tempFileData := []byte("change_test=1234")
			tempFilePath := filepath.Join(os.TempDir(), "temp.env")

			err := os.WriteFile(tempFilePath, tempFileData, 0644)
			defer os.Remove(tempFilePath)
			Expect(err).To(BeNil())

			// Create a configuration service
			svc := config.NewConfigService(config.ConfigOption{
				ConfigFilePathList: []string{tempFilePath},
				WatchConfigChange:  true,
			})

			// Check the initial value
			Expect(svc.GetInt("change_test")).To(Equal(1234))

			// Change the configuration file
			tempFileData = []byte("change_test=5678")
			err = os.WriteFile(tempFilePath, tempFileData, 0644)
			Expect(err).To(BeNil())

			// Check the changed value
			Eventually(func() int {
				return svc.GetInt("change_test")
			}).Should(Equal(5678))
		})
	})
})
