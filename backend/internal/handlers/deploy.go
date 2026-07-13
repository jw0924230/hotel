package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// DeployHandler handles frontend deployment trigger requests.
type DeployHandler struct{}

// NewDeployHandler creates a new DeployHandler.
func NewDeployHandler() *DeployHandler {
	return &DeployHandler{}
}

// TriggerDeploy triggers the GitHub Actions workflow to build and deploy the frontend.
// POST /api/deploy
func (h *DeployHandler) TriggerDeploy(c *fiber.Ctx) error {
	pat := os.Getenv("GH_PAT")
	owner := os.Getenv("GH_OWNER")
	repo := os.Getenv("GH_REPO")

	if pat == "" || owner == "" || repo == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "伺服器端未配置必要的環境變數 (GH_PAT, GH_OWNER, GH_REPO)",
		})
	}

	workflowFile := "frontend-deploy.yml"
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches", owner, repo, workflowFile)

	// GitHub Actions repository dispatch payload
	requestBody, err := json.Marshal(map[string]interface{}{
		"ref": "main",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "建立請求 JSON 失敗: " + err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "建立 GitHub 請求物件失敗: " + err.Error(),
		})
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+pat)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "呼叫 GitHub API 失敗: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("GitHub API 拒絕請求，狀態碼: %d (請確認 Token 權限及專案路徑)", resp.StatusCode),
		})
	}

	return c.JSON(fiber.Map{
		"message": "已成功觸發前端網站部署！請稍候 3~5 分鐘讓 GitHub Action 完成編譯。",
	})
}
