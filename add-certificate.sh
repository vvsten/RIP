#!/bin/bash
# Скрипт для добавления самоподписанного сертификата в Keychain macOS
# Это нужно для работы Tauri с HTTPS и самоподписанными сертификатами

CERT_FILE="certs/server.crt"
CERT_NAME="localhost (GruzDelivery)"

echo "Добавление сертификата $CERT_NAME в Keychain..."

# Удаляем старый сертификат, если есть
security delete-certificate -c "$CERT_NAME" login.keychain 2>/dev/null

# Добавляем сертификат в Keychain
security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain-db "$CERT_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Сертификат успешно добавлен в Keychain!"
    echo "Теперь Tauri сможет подключаться к HTTPS серверу с этим сертификатом."
else
    echo "❌ Ошибка при добавлении сертификата."
    echo "Попробуйте выполнить вручную:"
    echo "  security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain-db $CERT_FILE"
fi

