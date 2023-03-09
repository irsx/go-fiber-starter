package services

import (
	"encoding/json"
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/app/transformer"
	"go-fiber-starter/configs"
	"go-fiber-starter/constants"
	"go-fiber-starter/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
)

type NewsService struct{}

func (s *NewsService) newsRepo() *repository.NewsRepository {
	return new(repository.NewsRepository)
}

func (s *NewsService) redisStorage() *redis.Storage {
	return configs.RedisStorage
}

func (s *NewsService) List(ctx *fiber.Ctx, status int) (err error) {
	cacheKey := s.keyNewsListCache(status)
	cacheData := s.getNewsListCache(cacheKey)
	if cacheData != nil {
		return utils.JsonSuccess(ctx, cacheData)
	}

	news, err := s.newsRepo().GetAll(status)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_NEWS_LIST")
	}

	newsData := transformer.NewsListTransformer(news)
	if len(newsData) > 0 {
		s.setNewsListCache(cacheKey, newsData)
	}

	return utils.JsonSuccess(ctx, newsData)
}

func (s *NewsService) Add(ctx *fiber.Ctx, req dto.NewsRequestDTO) error {
	if err := req.Validate(); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	storeData := s.storeFromRequest(models.News{}, req)
	news, err := s.newsRepo().Insert(storeData)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_NEWS_ADD")
	}

	// RESET LIST CACHE
	s.resetNewsListCache()

	// SEND QUEUE
	utils.SendJobWithDefaultPayloads(constants.QueueNewsBroadcast, news)

	return utils.JsonSuccess(ctx, transformer.NewsDetailTransformer(news))
}

func (s *NewsService) Detail(ctx *fiber.Ctx, guid string) error {
	cacheKey := s.keyNewsDetailCache(guid)
	cacheData := s.getNewsDetailCache(cacheKey)
	if cacheData != nil {
		return utils.JsonSuccess(ctx, cacheData)
	}

	news, err := s.newsRepo().FindByGUID(guid)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, err)
	}

	newsDetail := transformer.NewsDetailTransformer(news)
	s.setNewsDetailCache(cacheKey, newsDetail)

	return utils.JsonSuccess(ctx, newsDetail)
}

func (s *NewsService) Update(ctx *fiber.Ctx, guid string, req dto.NewsRequestDTO) error {
	if err := req.Validate(); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	news, err := s.newsRepo().FindByGUID(guid)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, err)
	}

	storeData := s.storeFromRequest(news, req)
	updatedNews, err := s.newsRepo().UpdateByGUID(guid, storeData)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_NEWS_UPDATE")
	}

	// UPDATE CACHE
	newsDetail := transformer.NewsDetailTransformer(updatedNews)
	cacheKey := s.keyNewsDetailCache(guid)
	s.setNewsDetailCache(cacheKey, newsDetail)

	// RESET LIST CACHE
	s.resetNewsListCache()

	return utils.JsonSuccess(ctx, newsDetail)
}

func (s *NewsService) Delete(ctx *fiber.Ctx, guid string) error {
	_, err := s.newsRepo().UpdateByGUID(guid, models.News{Status: "0"})
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_NEWS_DELETE")
	}

	// DELETE CACHE
	cacheKey := s.keyNewsDetailCache(guid)
	s.redisStorage().Delete(cacheKey)

	// RESET LIST CACHE
	s.resetNewsListCache()

	return utils.JsonSuccess(ctx, fiber.Map{"guid": guid})
}

func (s *NewsService) storeFromRequest(news models.News, req dto.NewsRequestDTO) models.News {
	news.UserGUID = req.UserGUID
	news.Title = req.Title
	news.Description = req.Description
	news.Image = req.Image
	news.HyperLink = req.HyperLink
	news.Status = strconv.Itoa(req.Status)
	return news
}

func (s *NewsService) keyNewsListCache(status int) string {
	return constants.CacheNewsList + "_" + strconv.Itoa(status)
}

func (s *NewsService) keyNewsDetailCache(guid string) string {
	return constants.CacheNewsDetail + "_" + guid
}

func (s *NewsService) getNewsListCache(key string) (list []*transformer.NewsDetailResponse) {
	cacheNews, _ := s.redisStorage().Get(key)
	if len(cacheNews) > 0 {
		utils.Logger.Info("✅ GET CACHE NEWS LIST FROM KEY " + key)
		if err := json.Unmarshal(cacheNews, &list); err != nil {
			return nil
		}

		return list
	}

	return nil
}

func (s *NewsService) setNewsListCache(key string, news []transformer.NewsDetailResponse) {
	utils.Logger.Info("✅ SET CACHE NEWS LIST TO KEY " + key)
	cacheDataBytes, _ := json.Marshal(news)
	err := s.redisStorage().Set(key, cacheDataBytes, 12*time.Hour)
	if err != nil {
		utils.Logger.Error("❌ REDIS KEY " + key + " ERROR: " + err.Error())
	}
}

func (s *NewsService) getNewsDetailCache(key string) (detail *transformer.NewsDetailResponse) {
	cacheNews, _ := s.redisStorage().Get(key)
	if len(cacheNews) > 0 {
		utils.Logger.Info("✅ GET CACHE NEWS DETAIL FROM KEY " + key)
		if err := json.Unmarshal(cacheNews, &detail); err != nil {
			return nil
		}

		return detail
	}

	return nil
}

func (s *NewsService) setNewsDetailCache(key string, news *transformer.NewsDetailResponse) {
	utils.Logger.Info("✅ SET CACHE NEWS DETAIL TO KEY " + key)
	cacheDataBytes, _ := json.Marshal(news)
	err := s.redisStorage().Set(key, cacheDataBytes, 12*time.Hour)
	if err != nil {
		utils.Logger.Error("❌ REDIS KEY " + key + " ERROR: " + err.Error())
	}
}

func (s *NewsService) resetNewsListCache() {
	utils.Logger.Info("✅ RESET CACHE NEWS LIST")
	rs := s.redisStorage()
	rs.Delete(s.keyNewsListCache(0))
	rs.Delete(s.keyNewsListCache(1))
}
