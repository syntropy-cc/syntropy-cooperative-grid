/**
 * File: token_management.go
 * Purpose: Métodos de gerenciamento de Grid Token no SetupManager
 * Dependencies: token_manager.go, setup.go
 * Exports: Métodos públicos do SetupManager para gerenciamento de token
 * Author: Syntropy Setup Component
 * Created: 2025-01-27
 * Modified: 2025-01-27
 * Version: 1.0.0
 *
 * Business Context:
 * Este arquivo contém os métodos públicos do SetupManager para gerenciamento
 * de Grid Token, fornecendo uma interface simplificada para operações de
 * token através do CLI e outras interfaces.
 */

package setup

// Métodos de gerenciamento de Grid Token

// GenerateGridToken gera um novo Grid Token
func (sm *SetupManager) GenerateGridToken() (string, error) {
	sm.logger.LogStep("grid_token_generation_start", map[string]interface{}{
		"method": "manual_generation",
	})

	token, err := sm.tokenManager.GenerateToken()
	if err != nil {
		return "", sm.handleError(err, "grid_token_generation_failed")
	}

	if err := sm.tokenManager.SaveToken(token); err != nil {
		return "", sm.handleError(err, "grid_token_save_failed")
	}

	sm.logger.LogStep("grid_token_generation_completed", map[string]interface{}{
		"token_preview": token[:8] + "...[HIDDEN]",
	})

	return token, nil
}

// GetGridToken obtém o Grid Token atual
func (sm *SetupManager) GetGridToken() (string, error) {
	sm.logger.LogStep("grid_token_retrieval_start", nil)

	token, err := sm.tokenManager.LoadToken()
	if err != nil {
		return "", sm.handleError(err, "grid_token_load_failed")
	}

	sm.logger.LogStep("grid_token_retrieval_completed", map[string]interface{}{
		"token_preview": token[:8] + "...[HIDDEN]",
	})

	return token, nil
}

// RotateGridToken rotaciona o Grid Token atual
func (sm *SetupManager) RotateGridToken() (string, error) {
	sm.logger.LogStep("grid_token_rotation_start", nil)

	newToken, err := sm.tokenManager.RotateToken()
	if err != nil {
		return "", sm.handleError(err, "grid_token_rotation_failed")
	}

	sm.logger.LogStep("grid_token_rotation_completed", map[string]interface{}{
		"new_token_preview": newToken[:8] + "...[HIDDEN]",
	})

	return newToken, nil
}

// ExportGridToken exporta o Grid Token para backup
func (sm *SetupManager) ExportGridToken(outputPath string) error {
	sm.logger.LogStep("grid_token_export_start", map[string]interface{}{
		"output_path": outputPath,
	})

	if err := sm.tokenManager.ExportToken(outputPath); err != nil {
		return sm.handleError(err, "grid_token_export_failed")
	}

	sm.logger.LogStep("grid_token_export_completed", map[string]interface{}{
		"output_path": outputPath,
	})

	return nil
}

// ImportGridToken importa Grid Token de backup
func (sm *SetupManager) ImportGridToken(inputPath string) error {
	sm.logger.LogStep("grid_token_import_start", map[string]interface{}{
		"input_path": inputPath,
	})

	if err := sm.tokenManager.ImportToken(inputPath); err != nil {
		return sm.handleError(err, "grid_token_import_failed")
	}

	sm.logger.LogStep("grid_token_import_completed", map[string]interface{}{
		"input_path": inputPath,
	})

	return nil
}

// DeleteGridToken remove o Grid Token
func (sm *SetupManager) DeleteGridToken() error {
	sm.logger.LogStep("grid_token_deletion_start", nil)

	if err := sm.tokenManager.DeleteToken(); err != nil {
		return sm.handleError(err, "grid_token_deletion_failed")
	}

	sm.logger.LogStep("grid_token_deletion_completed", nil)

	return nil
}

// GridTokenExists verifica se existe um Grid Token
func (sm *SetupManager) GridTokenExists() (bool, error) {
	exists, err := sm.tokenManager.TokenExists()
	if err != nil {
		return false, sm.handleError(err, "grid_token_existence_check_failed")
	}

	return exists, nil
}
