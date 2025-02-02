package repository

import (
    "context"
    "encoding/json"
    "errors"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/0xPixelNinja/GinFusion/internal/config"
    "github.com/0xPixelNinja/GinFusion/internal/models"
)

var ctx = context.Background()
var redisClient *redis.Client

// InitRedis initializes the Redis client.
func InitRedis(cfg *config.Config) {
    redisClient = redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.Addr,
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })
}

// GetRedisClient returns the initialized Redis client.
func GetRedisClient() *redis.Client {
    return redisClient
}

// CreateUser stores a user in Redis and adds the user ID to a set.
func CreateUser(user models.User) error {
    if redisClient == nil {
        return errors.New("redis client not initialized")
    }
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }
    if err := redisClient.Set(ctx, "user:"+user.ID, data, 0).Err(); err != nil {
        return err
    }
    return redisClient.SAdd(ctx, "users", user.ID).Err()
}

// GetUser retrieves a user by ID.
func GetUser(userID string) (models.User, error) {
    var user models.User
    if redisClient == nil {
        return user, errors.New("redis client not initialized")
    }
    data, err := redisClient.Get(ctx, "user:"+userID).Result()
    if err != nil {
        return user, err
    }
    if err := json.Unmarshal([]byte(data), &user); err != nil {
        return user, err
    }
    return user, nil
}

// ListUsers returns all users.
func ListUsers() ([]models.User, error) {
    if redisClient == nil {
        return nil, errors.New("redis client not initialized")
    }
    userIDs, err := redisClient.SMembers(ctx, "users").Result()
    if err != nil {
        return nil, err
    }
    var users []models.User
    for _, id := range userIDs {
        user, err := GetUser(id)
        if err == nil {
            users = append(users, user)
        }
    }
    return users, nil
}

// UpdateUser updates an existing user's data.
func UpdateUser(user models.User) error {
    if redisClient == nil {
        return errors.New("redis client not initialized")
    }
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return redisClient.Set(ctx, "user:"+user.ID, data, 0).Err()
}

// GetUsageStats returns dummy usage statistics.
func GetUsageStats() (map[string]interface{}, error) {
    stats := map[string]interface{}{
        "total_users": 0,
    }
    if redisClient != nil {
        count, err := redisClient.SCard(ctx, "users").Result()
        if err == nil {
            stats["total_users"] = count
        }
    }
    return stats, nil
}

// CreateAPIKey stores an API key in Redis, adds it to the set, and maps the user to the key.
func CreateAPIKey(apiKey models.APIKey) error {
	if redisClient == nil {
		return errors.New("redis client not initialized")
	}
	data, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}
	if err := redisClient.Set(ctx, "apikey:"+apiKey.Key, data, 0).Err(); err != nil {
		return err
	}
	if err := redisClient.SAdd(ctx, "apikeys", apiKey.Key).Err(); err != nil {
		return err
	}
	return redisClient.Set(ctx, "user_apikey:"+apiKey.UserID, apiKey.Key, 0).Err()
}

// GetAPIKey retrieves an API key.
func GetAPIKey(key string) (models.APIKey, error) {
    var apiKey models.APIKey
    if redisClient == nil {
        return apiKey, errors.New("redis client not initialized")
    }
    data, err := redisClient.Get(ctx, "apikey:"+key).Result()
    if err != nil {
        return apiKey, err
    }
    if err := json.Unmarshal([]byte(data), &apiKey); err != nil {
        return apiKey, err
    }
    return apiKey, nil
}

// ListAPIKeys returns all API keys.
func ListAPIKeys() ([]models.APIKey, error) {
    if redisClient == nil {
        return nil, errors.New("redis client not initialized")
    }
    keys, err := redisClient.SMembers(ctx, "apikeys").Result()
    if err != nil {
        return nil, err
    }
    var apiKeys []models.APIKey
    for _, key := range keys {
        apiKey, err := GetAPIKey(key)
        if err == nil {
            apiKeys = append(apiKeys, apiKey)
        }
    }
    return apiKeys, nil
}

// UpdateAPIKey updates an existing API key.
func UpdateAPIKey(apiKey models.APIKey) error {
    if redisClient == nil {
        return errors.New("redis client not initialized")
    }
    data, err := json.Marshal(apiKey)
    if err != nil {
        return err
    }
    return redisClient.Set(ctx, "apikey:"+apiKey.Key, data, 0).Err()
}

// DeleteAPIKey removes an API key and its user mapping.
func DeleteAPIKey(key string) error {
	if redisClient == nil {
		return errors.New("redis client not initialized")
	}
	apiKey, err := GetAPIKey(key)
	if err != nil {
		return err
	}
	if err := redisClient.Del(ctx, "apikey:"+key).Err(); err != nil {
		return err
	}
	if err := redisClient.SRem(ctx, "apikeys", key).Err(); err != nil {
		return err
	}
	return redisClient.Del(ctx, "user_apikey:"+apiKey.UserID).Err()
}


// GetAPIKeyByUser retrieves the API key for a given user.
func GetAPIKeyByUser(userID string) (models.APIKey, error) {
	var apiKey models.APIKey
	key, err := redisClient.Get(ctx, "user_apikey:"+userID).Result()
	if err != nil {
		return apiKey, err
	}
	return GetAPIKey(key)
}



// AddTokenToBlacklist adds a JWT token to a blacklist in Redis with the specified TTL.
func AddTokenToBlacklist(token string, ttl time.Duration) error {
	return redisClient.Set(ctx, "blacklist:"+token, "blacklisted", ttl).Err()
}

// IsTokenBlacklisted checks if the given token is blacklisted.
func IsTokenBlacklisted(token string) bool {
	_, err := redisClient.Get(ctx, "blacklist:"+token).Result()
	return err == nil
}

// LogActivity logs a user's activity by pushing a new entry onto a Redis list.
func LogActivity(userID, activity string) error {
	key := "activity:" + userID
	entry := time.Now().Format(time.RFC3339) + " - " + activity
	return redisClient.LPush(ctx, key, entry).Err()
}

// GetActivityLog retrieves the recent activity log for a given user.
func GetActivityLog(userID string, count int64) ([]string, error) {
	key := "activity:" + userID
	return redisClient.LRange(ctx, key, 0, count-1).Result()
}
