package usb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// listDevicesLinux lista dispositivos USB no Linux
func listDevicesLinux() ([]USBDevice, error) {
	var devices []USBDevice

	// Usar lsblk para detectar dispositivos de forma mais robusta
	cmd := exec.Command("lsblk", "-J", "-o", "NAME,TYPE,SIZE,MODEL,VENDOR,SERIAL,RM,TRAN,HOTPLUG")
	output, err := cmd.Output()
	if err != nil {
		// Fallback para método antigo se lsblk falhar
		return listDevicesLinuxFallback()
	}

	var lsblkOutput struct {
		BlockDevices []struct {
			Name      string `json:"name"`
			Type      string `json:"type"`
			Size      string `json:"size"`
			Model     string `json:"model"`
			Vendor    string `json:"vendor"`
			Serial    string `json:"serial"`
			Removable bool   `json:"rm"`
			Transport string `json:"tran"`
			Hotplug   bool   `json:"hotplug"`
		} `json:"blockdevices"`
	}

	if err := json.Unmarshal(output, &lsblkOutput); err != nil {
		return listDevicesLinuxFallback()
	}

	for _, blockDev := range lsblkOutput.BlockDevices {
		// Filtrar apenas dispositivos de disco (não partições)
		if blockDev.Type != "disk" {
			continue
		}

		// Verificar se é removível ou USB
		isRemovable := blockDev.Removable ||
			strings.ToLower(blockDev.Transport) == "usb" ||
			blockDev.Hotplug

		// Incluir também dispositivos pequenos que podem ser USBs
		// (menos de 2TB e com transport USB ou removível)
		if !isRemovable {
			// Verificar tamanho para heurística adicional
			if sizeGB := parseSizeToGB(blockDev.Size); sizeGB > 0 && sizeGB < 2000 {
				// Verificar se tem características de USB
				if strings.Contains(strings.ToLower(blockDev.Model), "usb") ||
					strings.Contains(strings.ToLower(blockDev.Vendor), "usb") {
					isRemovable = true
				}
			}
		}

		if !isRemovable {
			continue
		}

		device := USBDevice{
			Path:      "/dev/" + blockDev.Name,
			Removable: true,
			Platform:  "linux",
			Model:     blockDev.Model,
			Vendor:    blockDev.Vendor,
			Serial:    blockDev.Serial,
			Size:      blockDev.Size,
			SizeGB:    parseSizeToGB(blockDev.Size),
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// listDevicesLinuxFallback método fallback para detecção de dispositivos
func listDevicesLinuxFallback() ([]USBDevice, error) {
	var devices []USBDevice

	// Listar dispositivos de bloco (sd*, nvme*, mmcblk*)
	patterns := []string{"/sys/block/sd*", "/sys/block/nvme*", "/sys/block/mmcblk*"}

	for _, pattern := range patterns {
		blockDevs, _ := filepath.Glob(pattern)
		for _, blockDev := range blockDevs {
			devName := filepath.Base(blockDev)

			// Verificar se é removível
			removableData, _ := os.ReadFile(filepath.Join(blockDev, "removable"))
			if strings.TrimSpace(string(removableData)) != "1" {
				// Para NVMe, verificar se é pequeno (heurística para USB)
				if strings.HasPrefix(devName, "nvme") {
					if sizeData, err := os.ReadFile(filepath.Join(blockDev, "size")); err == nil {
						if sectors, err := strconv.ParseInt(strings.TrimSpace(string(sizeData)), 10, 64); err == nil {
							sizeBytes := sectors * 512
							sizeGB := int(sizeBytes / (1024 * 1024 * 1024))
							// Se for menor que 2TB, pode ser USB
							if sizeGB >= 2000 {
								continue
							}
						}
					}
				} else {
					continue
				}
			}

			// Obter informações do dispositivo
			device := USBDevice{
				Path:      "/dev/" + devName,
				Removable: true,
				Platform:  "linux",
			}

			// Tamanho
			if sizeData, err := os.ReadFile(filepath.Join(blockDev, "size")); err == nil {
				if sectors, err := strconv.ParseInt(strings.TrimSpace(string(sizeData)), 10, 64); err == nil {
					sizeBytes := sectors * 512
					device.Size = formatSize(sizeBytes)
					device.SizeGB = int(sizeBytes / (1024 * 1024 * 1024))
				}
			}

			// Modelo e Vendor
			if modelData, err := os.ReadFile(filepath.Join(blockDev, "device/model")); err == nil {
				device.Model = strings.TrimSpace(string(modelData))
			}
			if vendorData, err := os.ReadFile(filepath.Join(blockDev, "device/vendor")); err == nil {
				device.Vendor = strings.TrimSpace(string(vendorData))
			}

			devices = append(devices, device)
		}
	}

	return devices, nil
}

// parseSizeToGB converte string de tamanho (ex: "1.5G") para GB
func parseSizeToGB(sizeStr string) int {
	if sizeStr == "" {
		return 0
	}

	// Remover espaços
	sizeStr = strings.TrimSpace(sizeStr)

	// Extrair número e unidade
	var size float64
	var unit string

	if strings.HasSuffix(sizeStr, "G") {
		fmt.Sscanf(sizeStr, "%f%s", &size, &unit)
		return int(size)
	} else if strings.HasSuffix(sizeStr, "M") {
		fmt.Sscanf(sizeStr, "%f%s", &size, &unit)
		return int(size / 1024)
	} else if strings.HasSuffix(sizeStr, "T") {
		fmt.Sscanf(sizeStr, "%f%s", &size, &unit)
		return int(size * 1024)
	}

	return 0
}

// createUSBWithNoCloudLinux cria USB usando estratégia NoCloud (dd + partição CIDATA)
func createUSBWithNoCloudLinux(devicePath string, config *Config, workDir, cacheDir string) error {
	// Verificar se tem permissões de root
	if os.Geteuid() != 0 {
		return fmt.Errorf("este comando requer privilégios de root (use sudo)")
	}

	// Gerenciar ISO com cache inteligente
	isoPath := config.ISOPath
	if isoPath == "" {
		var err error
		isoPath, err = manageISOCache(cacheDir)
		if err != nil {
			return fmt.Errorf("erro ao gerenciar ISO: %w", err)
		}
	}

	fmt.Println("📝 Gravando ISO Ubuntu no USB...")

	// Usar dd para gravar ISO no USB com timeout
	if err := runCommandWithTimeout(30*time.Minute, "dd",
		fmt.Sprintf("if=%s", isoPath),
		fmt.Sprintf("of=%s", devicePath),
		"bs=4M",
		"status=progress",
		"oflag=sync"); err != nil {
		return fmt.Errorf("erro ao gravar ISO: %w", err)
	}

	// Sync para garantir que tudo foi gravado
	exec.Command("sync").Run()

	fmt.Println("✅ ISO gravada com sucesso!")
	fmt.Println("📦 Criando partição CIDATA para cloud-init...")

	// Aguardar um pouco para o sistema reconhecer as mudanças
	time.Sleep(2 * time.Second)

	// Criar partição CIDATA usando sgdisk (GPT)
	if err := runCommandWithTimeout(30*time.Second, "sgdisk", "-e", devicePath); err != nil {
		return fmt.Errorf("erro ao reparar GPT: %w", err)
	}
	if err := runCommandWithTimeout(30*time.Second, "sgdisk", "-n", "0:0:+128MiB", "-t", "0:0700", "-c", "0:CIDATA", devicePath); err != nil {
		return fmt.Errorf("erro ao criar partição CIDATA: %w", err)
	}

	// Aguardar para a partição ser criada
	time.Sleep(1 * time.Second)

	// Determinar nome da partição CIDATA
	cidataPartition := devicePath + "2"
	if strings.Contains(devicePath, "nvme") {
		cidataPartition = devicePath + "p2"
	} else if strings.Contains(devicePath, "mmcblk") {
		cidataPartition = devicePath + "p2"
	}

	// Formatar partição CIDATA
	if err := runCommandWithTimeout(30*time.Second, "mkfs.vfat", "-F", "32", "-n", "CIDATA", cidataPartition); err != nil {
		return fmt.Errorf("erro ao formatar partição CIDATA: %w", err)
	}

	// Montar partição CIDATA
	mountPoint := filepath.Join(workDir, "cidata-mount")
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("erro ao criar ponto de montagem: %w", err)
	}

	if err := runCommandWithTimeout(15*time.Second, "mount", cidataPartition, mountPoint); err != nil {
		return fmt.Errorf("erro ao montar partição CIDATA: %w", err)
	}

	// Copiar arquivos cloud-init para a partição CIDATA
	fmt.Println("📝 Copiando arquivos cloud-init para partição CIDATA...")
	cloudInitDir := filepath.Join(workDir, "cloud-init")

	files := []string{"user-data", "meta-data", "network-config"}
	for _, file := range files {
		srcPath := filepath.Join(cloudInitDir, file)
		dstPath := filepath.Join(mountPoint, file)

		if err := runCommandWithTimeout(30*time.Second, "cp", srcPath, dstPath); err != nil {
			// Desmontar antes de retornar erro
			runCommandWithTimeout(10*time.Second, "umount", mountPoint)
			return fmt.Errorf("erro ao copiar %s: %w", file, err)
		}
	}

	// Desmontar partição CIDATA
	if err := runCommandWithTimeout(10*time.Second, "umount", mountPoint); err != nil {
		return fmt.Errorf("erro ao desmontar partição CIDATA: %w", err)
	}

	// Limpar ponto de montagem
	os.RemoveAll(mountPoint)

	// Sync final
	exec.Command("sync").Run()

	fmt.Println("✅ USB criado com sucesso usando estratégia NoCloud!")
	fmt.Println("🔧 O USB agora contém:")
	fmt.Println("   • ISO Ubuntu original (bootável)")
	fmt.Println("   • Partição CIDATA com configuração cloud-init")
	fmt.Println("   • Configuração será aplicada automaticamente no boot")

	return nil
}

// createUSBLinux função legada - redirecionar para nova implementação
func createUSBLinux(devicePath string, config *Config, workDir, cacheDir string) error {
	// Função legada - redirecionar para nova implementação
	return createUSBWithNoCloudLinux(devicePath, config, workDir, cacheDir)
}

// formatUSBLinux formata um dispositivo USB no Linux
func formatUSBLinux(devicePath, label string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("este comando requer privilégios de root (use sudo)")
	}

	// Desmontar todas as partições do dispositivo
	fmt.Println("🔓 Desmontando partições existentes...")
	umountAll(devicePath)

	// Aguardar um pouco para garantir que o dispositivo esteja livre
	time.Sleep(1 * time.Second)

	// Criar nova tabela de partições GPT (mais moderna)
	fmt.Println("📋 Criando tabela de partições...")
	if err := runCommandWithTimeout(30*time.Second, "parted", "-s", devicePath, "mklabel", "gpt"); err != nil {
		// Fallback para msdos se GPT falhar
		if err := runCommandWithTimeout(30*time.Second, "parted", "-s", devicePath, "mklabel", "msdos"); err != nil {
			return fmt.Errorf("erro ao criar tabela de partições: %w", err)
		}
	}

	// Criar partição primária com alinhamento adequado (1MiB)
	fmt.Println("📦 Criando partição...")
	if err := runCommandWithTimeout(30*time.Second, "parted", "-s", devicePath, "mkpart", "primary", "fat32", "1MiB", "100%"); err != nil {
		return fmt.Errorf("erro ao criar partição: %w", err)
	}

	// Aguardar para garantir que a partição seja criada
	time.Sleep(1 * time.Second)

	// Determinar nome da partição
	partition := devicePath + "1"
	if strings.Contains(devicePath, "nvme") {
		partition = devicePath + "p1"
	} else if strings.Contains(devicePath, "mmcblk") {
		partition = devicePath + "p1"
	}

	// Verificar se a partição existe
	if _, err := os.Stat(partition); err != nil {
		return fmt.Errorf("partição %s não foi criada: %w", partition, err)
	}

	// Formatar partição com FAT32
	fmt.Println("💾 Formatando partição...")
	if err := runCommandWithTimeout(30*time.Second, "mkfs.vfat", "-F", "32", "-n", label, partition); err != nil {
		return fmt.Errorf("erro ao formatar partição: %w", err)
	}

	// Sincronizar para garantir que tudo foi escrito
	exec.Command("sync").Run()

	fmt.Printf("✅ Dispositivo %s formatado com sucesso!\n", devicePath)
	return nil
}

// runCommandWithTimeout executa um comando com timeout
func runCommandWithTimeout(timeout time.Duration, name string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("comando %s excedeu o timeout de %v", name, timeout)
		}
		return fmt.Errorf("erro ao executar comando %s: %w", name, err)
	}

	return nil
}

// umountAll desmonta todas as partições de um dispositivo
func umountAll(device string) {
	// Desmonta qualquer partição montada desse device
	sh := fmt.Sprintf(`mount | awk '$1 ~ "^%s" {print $3}'`, device)
	out, _ := exec.Command("bash", "-lc", sh).Output()
	for _, mp := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if mp != "" {
			exec.Command("umount", mp).Run()
		}
	}
}
