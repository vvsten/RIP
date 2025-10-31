package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisService - сервис для работы с Redis
type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService - создание нового Redis сервиса
func NewRedisService(host string, port int, password string, db int) *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	return &RedisService{
		client: rdb,
		ctx:    context.Background(),
	}
}

// WriteJWTToBlacklist - добавление JWT токена в blacklist
func (r *RedisService) WriteJWTToBlacklist(token string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:jwt:%s", token)
	return r.client.Set(r.ctx, key, "1", expiration).Err()
}

// CheckJWTInBlacklist - проверка наличия JWT токена в blacklist
func (r *RedisService) CheckJWTInBlacklist(token string) (bool, error) {
	key := fmt.Sprintf("blacklist:jwt:%s", token)
	result := r.client.Get(r.ctx, key)
	if result.Err() == redis.Nil {
		return false, nil
	}
	if result.Err() != nil {
		return false, result.Err()
	}
	return true, nil
}

// StoreUserSession - сохранение сессии пользователя
func (r *RedisService) StoreUserSession(userUUID string, sessionData map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("session:user:%s", userUUID)
	return r.client.HMSet(r.ctx, key, sessionData).Err()
}

// GetUserSession - получение сессии пользователя
func (r *RedisService) GetUserSession(userUUID string) (map[string]string, error) {
	key := fmt.Sprintf("session:user:%s", userUUID)
	result := r.client.HGetAll(r.ctx, key)
	return result.Val(), result.Err()
}

// DeleteUserSession - удаление сессии пользователя
func (r *RedisService) DeleteUserSession(userUUID string) error {
	key := fmt.Sprintf("session:user:%s", userUUID)
	return r.client.Del(r.ctx, key).Err()
}

// StoreRefreshToken - сохранение refresh токена
func (r *RedisService) StoreRefreshToken(userUUID, refreshToken string, expiration time.Duration) error {
	key := fmt.Sprintf("refresh:user:%s", userUUID)
	return r.client.Set(r.ctx, key, refreshToken, expiration).Err()
}

// GetRefreshToken - получение refresh токена пользователя
func (r *RedisService) GetRefreshToken(userUUID string) (string, error) {
	key := fmt.Sprintf("refresh:user:%s", userUUID)
	result := r.client.Get(r.ctx, key)
	return result.Val(), result.Err()
}

// DeleteRefreshToken - удаление refresh токена
func (r *RedisService) DeleteRefreshToken(userUUID string) error {
	key := fmt.Sprintf("refresh:user:%s", userUUID)
	return r.client.Del(r.ctx, key).Err()
}

// Ping - проверка соединения с Redis
func (r *RedisService) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Close - закрытие соединения с Redis
func (r *RedisService) Close() error {
	return r.client.Close()
}

