# Explication du projet Golang

## Procédure d'exécution :

Pour installer les dépendances :
```bash
go mod download
```
Pour lancer le programme :
```bash
go run exam.go
```

Une fois que le programme est lancé, il lit le fichier `config.txt` qui contient les clés par défaut pour la lecture des fichiers, dossiers, sortie de log et extension.
Vous pouvez choisir une des options disponibles dans le menu (Choix A, Choix B, Choix C, Choix D) ou Quitter pour fermer le programme.

## Fonctionnalités implémentées :

- **Choix A** : Analyse d'un fichier et affiche des infos / Stats de mots / Compteur de lignes / Filtres mot clé / Head et Tail
- **Choix B** : Analyse de tout les .txt dans un dossier / Rapport global / Indéxation / Fusion de tout les .txt
- **Choix C** : Génération d'une page Wikipédia et afficher les stats
- **Choix D** : Lister les processus / Recherche et filtrer les processus / Kill sécurisé

Niveau visé 14/20 

Description du travail effectué : 

### readFile
Ouvre le fichier avec `os.Open`, parcourt ligne par ligne avec `bufio.Scanner`, ignore les lignes vides et celles commençant par `#`, puis renvoie le contenu concaténé avec `strings.Join`.

### readConfig / defaultValues
`readConfig` : découpe le texte par `\n`, pour chaque ligne non vide et non commentée, utilise `SplitN` sur `=` pour extraire clé/valeur et remplit une map. `defaultValues` : complète les clés manquantes avec des valeurs par défaut.

### linesCount
Lit le fichier via `readFile`, découpe par `\n`. Si la dernière ligne est vide (retour à la ligne final), soustrait 1 au nombre de lignes pour éviter de compter une ligne fictive.

### totalWordsWithoutNumbers
Utilise `strings.Fields` pour découper le contenu en mots. Pour chaque mot, parcourt les caractères avec `unicode.IsDigit` : si toutes sont des chiffres, le mot est ignoré. Compte uniquement les mots contenant au moins un caractère non numérique.

### countLinesWithKeyword
Découpe le contenu en lignes, parcourt chaque ligne et incrémente un compteur si `strings.Contains(line, keyword)` est vrai.

### filterLinesWithKeyword / filterLinesWithoutKeyword
Créent le dossier de sortie avec `os.Mkdir`, lisent le fichier, découpent en lignes. Parcourent les lignes et écrivent dans le fichier de sortie celles qui contiennent (ou ne contiennent pas) le mot clé via `strings.Contains`. Utilisation de `defer file.Close()` pour fermer le fichier.

### headLines / tailLines
`headLines` : découpe en lignes puis prend le slice `lines[:n]`. `tailLines` : prend le slice `lines[len(lines)-n:]`. Écrivent chaque ligne dans le fichier de sortie.

### analyseAllTxtFiles
`os.ReadDir` pour lister le dossier. Boucle sur les entrées : ignore les sous-dossiers, filtre avec `strings.HasSuffix` pour ne garder que les `.txt`, récupère les infos avec `os.Lstat` et affiche nom + taille.

### reportGlobalFolder
`os.Stat` pour les infos du dossier, `os.ReadDir` pour les entrées. Boucle pour sommer les tailles et compter les fichiers (sous-dossiers exclus). Crée `report.txt` avec `fmt.Sprintf` pour écrire le rapport.

### listFiles
Parcourt les entrées du dossier avec `ReadDir`, affiche chaque fichier (nom, taille en Ko, date). Réécrit les mêmes infos dans `index.txt` avec `WriteString`.

### mergeFiles
Lit le dossier avec `ReadDir`, crée `merged.txt`. Parcourt tous les fichiers, appelle `readFile` sur chacun et concatène le contenu dans le fichier de sortie (avec `\n` entre chaque fichier).

### analysePageWikipedia
Construit l'URL avec `url.PathEscape`. Crée une requête HTTP GET avec `User-Agent`. Exécute la requête, vérifie le code 200. Parse le HTML avec `goquery.NewDocumentFromReader`. `doc.Find("p")` pour sélectionner les paragraphes, `.Each` pour extraire le texte. Filtre les paragraphes vides, joint avec `\n\n`, écrit dans le fichier. Réutilise `linesCount` et `totalWordsWithoutNumbers` pour les stats.

### listProcesses
Détecte l'OS avec `runtime.GOOS`. Lance `tasklist` (Windows) ou `ps -Ao pid,comm` (macOS) via `exec.Command` et `Output()`. Normalise les retours à la ligne (`\r\n` → `\n`) et découpe en lignes.

### searchProcess
Selon l'OS : Windows utilise `tasklist /FI IMAGENAME eq <process>`, macOS utilise `ps -Ao pid,comm`. Affiche la sortie brute.

### killProcess
Vérifie que le PID > 0. Demande une confirmation via `bufio.Scanner`. Si l'utilisateur tape `o`, lance `taskkill /F /PID` (Windows) ou `kill -9` (macOS). Analyse la sortie avec `errors.Is` et `strings.Contains` pour distinguer les erreurs (accès refusé, processus introuvable).
