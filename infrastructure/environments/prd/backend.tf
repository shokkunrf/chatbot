terraform {
  backend "gcs" {
    prefix = "chatbot-prd"
  }
}
