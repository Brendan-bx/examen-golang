package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	data, err := os.ReadFile("config.txt")
	if err != nil {
		log.Fatal(err)
	}

	config := defaultValues(readConfig(string(data)))
	readFile("config.txt")

	for {
		choix := menu(config)

		if choix == "Quitter" || strings.EqualFold(choix, "Quitter") {
			fmt.Println("Au revoir !")
			break
		}
	}
}

func readConfig(data string) map[string]string {
	lines := strings.Split(data, "\n")
	config := make(map[string]string)
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "=")
		config[parts[0]] = parts[1]
	}
	return config
}

func readFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return strings.Join(lines, "\n")
}

func defaultValues(config map[string]string) map[string]string {
	if config["default_file"] == "" {
		config["default_file"] = "data/input.txt"
	}
	if config["base_dir"] == "" {
		config["base_dir"] = "data"
	}
	if config["out_dir"] == "" {
		config["out_dir"] = "out"
	}
	if config["default_ext"] == "" {
		config["default_ext"] = ".txt"
	}
	return config
}

func menu(config map[string]string) string {
	choixMenu := bufio.NewScanner(os.Stdin)
	fmt.Println("=== Menu ===")
	fmt.Println("Choix A - Analyse d'un fichier")
	fmt.Println("Choix B - Analyse multi fichiers")
	fmt.Println("Choix C - Analyse d'une page wikipédia")
	fmt.Println("Choix D - Sous-menu ProcessOps")
	fmt.Println("Quitter")
	fmt.Print("================================================ \n")

	choixMenu.Scan()
	switch choixMenu.Text() {
	case "Choix A":
		fmt.Println("Veuillez choisir un fichier à analyser")
		choixFile := bufio.NewScanner(os.Stdin)
		choixFile.Scan()

		fileInfo, err := os.Lstat(choixFile.Text())
		defaultFile := config["default_file"]
		outDir := config["out_dir"]
		filePath := choixFile.Text()
		if filePath == "" {
			filePath = defaultFile
		}
		fileInfo, err = os.Lstat(filePath)
		if err != nil {
			fmt.Println("File not found")
			break
		}
		fmt.Println("Fichier trouvé: ", fileInfo.Name())
		fmt.Println("================================================ \n")
		choixOption := bufio.NewScanner(os.Stdin)
		fmt.Println("Veuillez choisir une option à analyser")
		fmt.Println("1. Informations du fichier")
		fmt.Println("2. Stats mots")
		fmt.Println("3. Compter lignes avec un mot clé")
		fmt.Println("4. Filtrer les lignes avec un mot clé")
		fmt.Println("5. Filtrer les lignes sans un mot clé")
		fmt.Println("6. Afficher les n premières lignes")
		fmt.Println("7. Afficher les n dernières lignes")
		fmt.Print("================================================ \n")

		choixOption.Scan()
		switch choixOption.Text() {
		case "1":
			fmt.Println("Taille du fichier: ", fileInfo.Size(), "Ko")
			fmt.Println("Date de modification: ", fileInfo.ModTime().Format("02/01/2006 15:04:05"))
			fmt.Println("Nombre de lignes: ", linesCount(filePath))
		case "2":
			fmt.Println("Nombre de mots: ", totalWordsWithoutNumbers(filePath))
			fmt.Println("Moyenne de mots par ligne: ", totalWordsWithoutNumbers(filePath)/linesCount(filePath))
		case "3":
			fmt.Println("Veuillez entrer un mot clé à rechercher")
			choixMot := bufio.NewScanner(os.Stdin)
			choixMot.Scan()
			fmt.Println("Nombre de lignes avec le mot clé: ", countLinesWithKeyword(filePath, choixMot.Text()))
		case "4":
			fmt.Println("Veuillez entrer un mot clé à rechercher")
			choixMot := bufio.NewScanner(os.Stdin)
			choixMot.Scan()
			filterLinesWithKeyword(filePath, choixMot.Text(), outDir)
			fmt.Println("Fichier filtré avec succès")
		case "5":
			fmt.Println("Veuillez entrer un mot clé à rechercher")
			choixMot := bufio.NewScanner(os.Stdin)
			choixMot.Scan()
			filterLinesWithoutKeyword(filePath, choixMot.Text(), outDir)
			fmt.Println("Fichier filtré avec succès")
		case "6":
			fmt.Println("Veuillez entrer le nombre de lignes à afficher")
			choixN := bufio.NewScanner(os.Stdin)
			choixN.Scan()
			n, err := strconv.Atoi(choixN.Text())
			if err != nil {
				fmt.Println("Veuillez entrer un nombre entier valide")
				break
			}
			headLines(filePath, n, outDir)
			fmt.Println("Fichier affiché avec succès")
		case "7":
			fmt.Println("Veuillez entrer le nombre de lignes à afficher")
			choixN := bufio.NewScanner(os.Stdin)
			choixN.Scan()
			n, err := strconv.Atoi(choixN.Text())
			if err != nil {
				fmt.Println("Veuillez entrer un nombre entier valide")
				break
			}
			tailLines(filePath, n, outDir)
			fmt.Println("Fichier affiché avec succès")
		}
	case "Choix B":
		fmt.Println("Veuillez entrer le nom du dossier à analyser")
		choixFolder := bufio.NewScanner(os.Stdin)
		choixFolder.Scan()
		defaultFolder := config["base_dir"]
		outDir := config["out_dir"]
		folderPath := choixFolder.Text()
		if folderPath == "" {
			folderPath = defaultFolder
		}
		fmt.Println("Dossier trouvé: ", folderPath)
		fmt.Println("================================================ \n")
		choixOption := bufio.NewScanner(os.Stdin)
		fmt.Println("Veuillez choisir une option à analyser")
		fmt.Println("1. Analyser tout les txt du dossier")
		fmt.Println("2. Rapport global du dossier")
		fmt.Println("3. Lister les fichiers du dossier")
		fmt.Println("4. Fusionner les fichiers txt du dossier")
		fmt.Print("================================================ \n")
		choixOption.Scan()
		switch choixOption.Text() {
		case "1":
			analyseAllTxtFiles(folderPath)
		case "2":
			reportGlobalFolder(folderPath)
		case "3":
			listFiles(folderPath, outDir)
		case "4":
			mergeFiles(folderPath, outDir)
		default:
			fmt.Println("Choisir un choix valide")

		}
	case "Choix C":
		fmt.Println("Veuillez entrer le(s) nom(s) de page(s) wikipédia à analyser")
		fmt.Println("(séparés par des virgules)")
		choixPage := bufio.NewScanner(os.Stdin)
		choixPage.Scan()
		raw := choixPage.Text()
		outDir := config["out_dir"]
		if raw == "" {
			fmt.Println("Veuillez entrer au moins un nom de page wikipédia valide")
			break
		}
		pages := strings.Split(raw, ",")
		for _, p := range pages {
			page := strings.TrimSpace(p)
			if page == "" {
				continue
			}
			analysePageWikipedia(page, outDir)
		}
	case "Choix D":
		fmt.Println("Veuillez choisir une option à analyser")
		fmt.Println("1. Lister les processus en cours")
		fmt.Println("2. Chercher un processus")
		fmt.Println("3. Tuer un processus")
		fmt.Println("4. Quitter")
		choixProcessOps := bufio.NewScanner(os.Stdin)
		choixProcessOps.Scan()
		switch choixProcessOps.Text() {
		case "1":
			fmt.Println("Veuillez entrer le nombre de processus à afficher")
			choixTopN := bufio.NewScanner(os.Stdin)
			choixTopN.Scan()
			topN, err := strconv.Atoi(choixTopN.Text())
			if err != nil {
				fmt.Println("Veuillez entrer un nombre entier valide")
			}
			listProcesses(topN)
		case "2":
			fmt.Println("Veuillez entrer le nom du processus à afficher")
			choixProcess := bufio.NewScanner(os.Stdin)
			choixProcess.Scan()
			process := choixProcess.Text()
			searchProcess(process)
		case "3":
			fmt.Println("Veuillez entrer le pid du processus à tuer")
			choixProcess := bufio.NewScanner(os.Stdin)
			choixProcess.Scan()
			process := choixProcess.Text()
			pid, err := strconv.Atoi(process)
			if err != nil {
				fmt.Println("Veuillez entrer un nombre entier valide")
			}
			killProcess(pid)
		default:
			fmt.Println("Choisir un choix valide")
		}
	}

	return choixMenu.Text()
}

func linesCount(filePath string) int {
	content := readFile(filePath)
	if len(content) == 0 {
		return 0
	}
	lines := strings.Split(content, "\n")

	if lines[len(lines)-1] == "" {
		return len(lines) - 1
	}
	return len(lines)
}

func totalWordsWithoutNumbers(filePath string) int {
	content := readFile(filePath)
	if len(content) == 0 {
		return 0
	}
	words := strings.Fields(content)
	count := 0
	for _, w := range words {
		isNumeric := true
		for _, r := range w {
			if !unicode.IsDigit(r) {
				isNumeric = false
				break
			}
		}
		if !isNumeric {
			count++
		}
	}
	return count
}

func countLinesWithKeyword(filePath string, keyword string) int {
	content := readFile(filePath)
	if len(content) == 0 {
		return 0
	}
	lines := strings.Split(content, "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, keyword) {
			count++
		}
	}
	return count
}

func filterLinesWithKeyword(filePath string, keyword string, outDir string) {
	os.Mkdir(outDir, 0755)
	content := readFile(filePath)
	fileName := "filtered.txt"
	if len(content) == 0 {
		return
	}
	lines := strings.Split(content, "\n")
	file, err := os.Create(outDir + "/" + fileName)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}
	defer file.Close()
	for _, line := range lines {
		if strings.Contains(line, keyword) {
			_, err := file.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Erreur lors de l'écriture dans le fichier")
				return
			}
		}
	}
}

func filterLinesWithoutKeyword(filePath string, keyword string, outDir string) {
	os.Mkdir(outDir, 0755)
	content := readFile(filePath)
	fileName := "filtered_not.txt"
	if len(content) == 0 {
		return
	}
	lines := strings.Split(content, "\n")
	file, err := os.Create(outDir + "/" + fileName)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}
	defer file.Close()
	for _, line := range lines {
		if !strings.Contains(line, keyword) {
			_, err := file.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Erreur lors de l'écriture dans le fichier")
				return
			}
		}
	}
}

func headLines(filePath string, n int, outDir string) {
	content := readFile(filePath)
	lines := strings.Split(content, "\n")
	file, err := os.Create(outDir + "/" + "head.txt")
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}
	defer file.Close()
	for _, line := range lines[:n] {
		file.WriteString(line + "\n")
	}
}

func tailLines(filePath string, n int, outDir string) {
	content := readFile(filePath)
	lines := strings.Split(content, "\n")
	file, err := os.Create(outDir + "/" + "tail.txt")
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}
	defer file.Close()
	for _, line := range lines[len(lines)-n:] {
		file.WriteString(line + "\n")
	}
}

func analyseAllTxtFiles(folderPath string) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier")
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		fileInfo, err := os.Lstat(folderPath + "/" + file.Name())
		if err != nil {
			fmt.Println("Erreur lors de la lecture du fichier")
			continue
		}
		fmt.Println(file.Name(), " - ", fileInfo.Size(), "Ko")
	}
}

func reportGlobalFolder(folderPath string) {
	dirInfo, err := os.Stat(folderPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier:", err)
		return
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du contenu du dossier:", err)
		return
	}

	var totalSize int64
	fileCount := 0

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fileCount++
		info, err := e.Info()
		if err != nil {
			continue
		}
		totalSize += info.Size()
	}

	name := dirInfo.Name()
	if name == "" || name == "." {
		name = filepath.Base(folderPath)
	}

	fmt.Println("=== Rapport global du dossier ===")
	fmt.Println("Nom du dossier   :", name)
	fmt.Println("Chemin           :", folderPath)
	fmt.Println("Taille totale    :", totalSize, "octets")
	if fileCount > 0 {
		fmt.Println("Taille moyenne   :", totalSize/int64(fileCount), "octets par fichier")
	}
	fmt.Println("Nombre de fichiers :", fileCount)
	fmt.Println("Droits (mode)    :", dirInfo.Mode().String())
	fmt.Println("Date modification:", dirInfo.ModTime().Format("02/01/2006 15:04:05"))
	fmt.Println("=================================")
}

func listFiles(folderPath string, outDir string) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier")
		return
	}

	for _, entry := range files {
		info, err := entry.Info()
		if err != nil {
			fmt.Println(entry.Name(), "- Erreur lors de l'obtention des infos du fichier")
			continue
		}
		fmt.Println(entry.Name(), " - ", info.Size()/1024, "Ko", " - ", info.ModTime().Format("02/01/2006 15:04:05"))
	}

	outFile, err := os.Create(outDir + "/" + "index.txt")
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}
	defer outFile.Close()
	for _, entry := range files {
		info, err := entry.Info()
		if err != nil {
			fmt.Println(entry.Name(), "- Erreur lors de l'obtention des infos du fichier")
			continue
		}
		outFile.WriteString(entry.Name() + " - " + strconv.FormatInt(info.Size()/1024, 10) + "Ko - " + info.ModTime().Format("02/01/2006 15:04:05") + "\n")
	}
	fmt.Println("Index.txt créé avec succès")
}

func mergeFiles(folderPath string, outDir string) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier")
		return
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
	}
	outFile, err := os.Create(outDir + "/" + "merged.txt")
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}
	defer outFile.Close()
	for _, file := range files {
		content := readFile(folderPath + "/" + file.Name())
		outFile.WriteString(content + "\n")
	}
	fmt.Println("Fichiers fusionnés avec succès")
}

func analysePageWikipedia(page string, outDir string) {
	wikiURL := "https://fr.wikipedia.org/wiki/" + url.PathEscape(page)

	req, err := http.NewRequest(http.MethodGet, wikiURL, nil)
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête HTTP Wikipédia :", err)
		return
	}
	req.Header.Set("User-Agent", "examen-go/1.0 (https://example.com)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erreur HTTP lors du téléchargement de la page Wikipédia :", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Erreur HTTP lors de la lecture de la page Wikipédia :", resp.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors du parsing HTML de la page Wikipédia :", err)
		return
	}

	var paragraphs []string
	doc.Find("p").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			paragraphs = append(paragraphs, text)
		}
	})

	if len(paragraphs) == 0 {
		fmt.Println("Aucun paragraphe texte trouvé sur la page Wikipédia :", page)
		return
	}

	_ = os.Mkdir(outDir, 0755)
	outPath := outDir + "/" + "wiki_" + page + ".txt"

	outFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier :", err)
		return
	}
	defer outFile.Close()

	if _, err := outFile.WriteString(strings.Join(paragraphs, "\n\n")); err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier Wikipédia :", err)
		return
	}

	nbLignes := linesCount(outPath)
	nbMots := totalWordsWithoutNumbers(outPath)

	fmt.Println("Texte Wikipédia extrait avec succès dans", outPath)
	fmt.Println(" - Nombre de lignes :", nbLignes)
	fmt.Println(" - Nombre de mots (hors nombres) :", nbMots)
}

func listProcesses(topN int) {
	gos := runtime.GOOS
	var output []byte
	var err error
	switch gos {
	case "windows":
		cmd := exec.Command("tasklist")
		output, err = cmd.Output()
	case "darwin":
		cmd := exec.Command("ps", "-Ao", "pid,comm")
		output, err = cmd.Output()
	}
	if err != nil {
		fmt.Println("Erreur lors de la lecture des processus :", err)
		return
	}
	text := strings.ReplaceAll(strings.TrimSuffix(string(output), "\n"), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	max := topN + 3
	if max > len(lines) {
		max = len(lines)
	}
	for _, line := range lines[:max] {
		fmt.Println(strings.TrimRight(line, "\r"))
	}
}

func searchProcess(process string) {
	gos := runtime.GOOS
	switch gos {
	case "windows":
		cmd := exec.Command("tasklist", "/FI", "STATUS eq running", "/FI", "IMAGENAME eq "+process)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Erreur lors de la recherche du processus :", err)
			return
		}
		fmt.Println(string(output))
	case "darwin":
		cmd := exec.Command("ps", "-Ao", "pid,comm | grep", process)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Erreur lors de la recherche du processus :", err)
			return
		}
		fmt.Println(string(output))
	default:
		fmt.Println("Système non supporté pour la recherche de processus")
		return
	}
}

func killProcess(pid int) {
	if pid <= 0 {
		fmt.Println("PID invalide :", pid)
		return
	}

	fmt.Println("Voulez-vous vraiment tuer ce processus ? (o/n)")
	choix := bufio.NewScanner(os.Stdin)
	choix.Scan()
	if choix.Text() != "o" {
		fmt.Println("Processus non tué")
		return
	}
	gos := runtime.GOOS
	pidStr := strconv.Itoa(pid)

	var cmd *exec.Cmd
	switch gos {
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/PID", pidStr)
	case "darwin":
		cmd = exec.Command("kill", "-9", pidStr)
	default:
		fmt.Println("Commande de kill non disponible pour ce système")
		return
	}

	output, err := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(output))

	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Println("Commande non disponible sur ce système")
			return
		}

		if strings.Contains(outStr, "Access is denied") ||
			strings.Contains(outStr, "Operation not permitted") {
			fmt.Println("Droits insuffisants pour tuer le processus", pid)
			return
		}

		if strings.Contains(outStr, "does not exist") ||
			strings.Contains(outStr, "No such process") {
			fmt.Println("Le processus", pid, "est déjà terminé ou introuvable")
			return
		}

		fmt.Println("Erreur lors de la tentative de kill du processus :", outStr)
		return
	}

	fmt.Println("Processus", pid, "tué avec succès")
}
