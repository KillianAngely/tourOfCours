Parfait, voilà tes deux énoncés. Format : objectif clair, étapes progressives, **pas de code**, pas de solution. À toi de te débrouiller.

---

# 🟢 Exercice intermédiaire — Chat Room WebSocket

## Contexte

Tu construis un serveur de chat où plusieurs clients se connectent via WebSocket, et chaque message envoyé par un client est diffusé à **tous les autres clients connectés**. C'est le pattern fan-out 1→N, en pur Go, sans WebRTC.

## Stack imposée

- `net/http` (stdlib) pour le serveur HTTP.
- `github.com/coder/websocket` ou `github.com/gorilla/websocket` pour les WebSockets (au choix, coder/websocket est plus moderne).
- Stdlib pour tout le reste : `sync`, `context`, channels.

## Critères de réussite

À la fin, tu dois pouvoir :

1. Lancer ton serveur Go.
2. Ouvrir plusieurs onglets de navigateur (ou plusieurs clients `websocat`) qui se connectent à `ws://localhost:8080/ws`.
3. Quand un client envoie un message, **tous les autres clients connectés le reçoivent** en temps réel.
4. Quand un client ferme son onglet, le serveur le détecte et le retire proprement — pas de fuite de goroutine, pas de panic.

## Les 5 étapes

### Étape 1 — Serveur WebSocket basique (echo)

- Un endpoint `/ws` qui upgrade la connexion HTTP en WebSocket.
- Pour chaque client connecté, lance une goroutine qui lit ses messages et lui renvoie **le même message** (echo).
- Pas encore de broadcast. Juste : ça marche, je reçois ce que j'envoie.
- Teste avec `websocat ws://localhost:8080/ws` depuis 2 terminaux séparés.

### Étape 2 — Le Hub central

Maintenant tu veux que les messages aillent **vers tous les clients**, pas juste vers leur émetteur.

- Crée une struct `Hub` qui maintient la liste des clients connectés.
- Quand un client se connecte, il s'enregistre auprès du hub.
- Quand il se déconnecte, il se désenregistre.
- Pour l'instant, l'enregistrement peut juste être une `map[*Client]bool` protégée par une `sync.Mutex`.

**Question à te poser** : où vit le Hub ? Combien y en a-t-il ? (Indice : un seul, partagé par tout le serveur.)

### Étape 3 — Le broadcast (fan-out)

- Chaque `Client` possède son **propre channel** pour recevoir les messages à envoyer (ex: `send chan []byte`).
- Chaque client a **2 goroutines** :
  - Une qui **lit** depuis la WebSocket et pousse les messages reçus vers le hub.
  - Une qui **écrit** vers la WebSocket en lisant depuis son channel `send`.
- Quand le hub reçoit un message d'un client, il fait une boucle sur tous les clients enregistrés et pousse le message dans leur channel `send`.

**Questions à te poser** :

- Pourquoi 2 goroutines par client et pas une seule ? (Indice : `Read` et `Write` sur une WebSocket peuvent être appelés en parallèle, mais pas la même opération en parallèle.)
- Channel `send` bufferisé ou pas ? Quelle taille ?
- Si un client est lent à consommer son channel `send`, qu'est-ce qui se passe pour le hub qui broadcast ? Comment éviter de bloquer ?

### Étape 4 — Vire la Mutex, passe en mode "tout-channels"

C'est ici que tu vas vraiment apprendre l'idiome Go avancé : remplacer les locks par des channels.

- Au lieu d'une `map[*Client]bool` + mutex, fais une **goroutine "hub"** qui possède la map et qui est la **seule à y toucher**.
- Cette goroutine boucle sur un `select` qui écoute 3 channels :
  - `register chan *Client` — un nouveau client veut s'enregistrer.
  - `unregister chan *Client` — un client part.
  - `broadcast chan []byte` — un message à diffuser.
- Personne d'autre ne touche à la map. Tout passe par ces channels.

C'est le pattern **"share memory by communicating"** (la phrase culte de Go). Très important à maîtriser.

### Étape 5 — Déconnexions propres avec `context`

- Quand un client ferme sa WebSocket (ou son réseau coupe), tu dois :
  - Détecter la déconnexion (lecture qui retourne une erreur).
  - Envoyer le client au channel `unregister` du hub.
  - Arrêter les 2 goroutines de ce client proprement (pas de leak).
  - Fermer son channel `send`.
- Utilise un `context.Context` par client pour orchestrer l'arrêt des goroutines de lecture/écriture.

**Test à faire** : lance ton serveur, connecte 3 clients, kill un client brutalement (`Ctrl+C` dans `websocat`). Vérifie avec `runtime.NumGoroutine()` toutes les 5 secondes que le compte redescend bien. Si ça grimpe indéfiniment, t'as une fuite.

## Bonus

- Une route `/stats` qui renvoie le nombre de clients connectés en JSON.
- Des "rooms" : chaque client envoie d'abord son nom de room, et ne reçoit que les messages de sa room (`map[string]map[*Client]bool`).
- Une page HTML simple servie par ton serveur Go qui te permet de tester depuis le navigateur sans avoir besoin de `websocat`.

## Quand tu auras fini

Tu dois pouvoir expliquer **avec tes mots** :

- Pourquoi un channel par client et pas un channel global.
- Pourquoi 2 goroutines par client.
- Pourquoi le pattern "goroutine hub" est mieux qu'un mutex.
- Comment tu détectes une déconnexion et comment tu nettoies.

Si l'un de ces 4 points te paraît flou, refais l'étape concernée.
