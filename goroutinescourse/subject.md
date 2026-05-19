## Projet : Agrégateur de news multi-sources

### Objectif fonctionnel

Tu construis un outil qui, étant donné un mot-clé, va chercher des articles/posts sur **plusieurs sources publiques en parallèle**, agrège tout et retourne un résultat dédoublonné et trié par date.

### Sources suggérées (toutes gratuites, sans clé API)

- **Hacker News** (Algolia Search API) : `https://hn.algolia.com/api/v1/search?query=MOTCLE`
- **Reddit** : `https://www.reddit.com/search.json?q=MOTCLE`
- **Dev.to** : `https://dev.to/api/articles?tag=MOTCLE`
- **GitHub** (repos) : `https://api.github.com/search/repositories?q=MOTCLE`

Tu peux en ajouter d'autres si tu veux, ou en retirer.

### Le défi en 5 étapes

À chaque étape, tu fais marcher le truc, tu mesures le temps d'exécution, puis tu passes à la suivante.

**Étape 1 — Version séquentielle**

- Une fonction `Search(keyword)` qui appelle les 4 sources **l'une après l'autre**
- Retourne une liste agrégée d'articles
- Pas de goroutines. Juste du Go basique.
- Mesure le temps total.

**Étape 2 — Parallélisation avec WaitGroup**

- Lance les 4 appels **en parallèle** avec des goroutines
- Synchronise avec `sync.WaitGroup`
- Protège l'accès à la liste partagée avec `sync.Mutex`
- Compare le temps avec l'étape 1.
- **Piège à éviter** : la closure dans la boucle (cherche "go loop variable capture")

**Étape 3 — Refactor avec channels**

- Vire la mutex
- Chaque goroutine envoie ses résultats dans un **channel**
- La fonction principale les lit et les agrège
- Réfléchis : channel bufferé ou pas ?
- Comment savoir quand toutes les goroutines ont fini ?

**Étape 4 — Context (timeout + annulation)**

- La fonction `Search` accepte maintenant un `context.Context` en premier paramètre
- Ajoute un timeout global de 2 secondes
- Si une source ne répond pas à temps, on l'ignore et on retourne ce qu'on a
- Utilise `select` pour écouter `ctx.Done()`
- Propage le `ctx` jusqu'aux appels HTTP (cherche `http.NewRequestWithContext`)

**Étape 5 — Worker pool + streaming**

- Limite à **2 requêtes HTTP simultanées max** (sémaphore avec channel)
- Au lieu de retourner un slice à la fin, **stream** les articles un par un dans un channel de sortie
- L'appelant peut afficher/traiter les résultats au fur et à mesure
- Ferme proprement le channel quand tout est fini

### Bonus (si tu veux pousser)

- Wrap le tout dans un handler HTTP `/search?q=...`
- Utilise `r.Context()` du handler pour propager l'annulation
- Teste avec `curl --max-time 1 http://localhost:8080/search?q=go` et vérifie que tes goroutines s'arrêtent bien (mets des `log.Println` pour le voir)
- Implémente du Server-Sent Events (SSE) pour streamer les résultats au client en temps réel

### Conseils pour bien apprendre

- À chaque étape, **mesure** avec `time.Since(start)`. Le déclic vient quand tu vois 3s → 1s en passant à l'étape 2.
- Mets des `log.Println` partout (genre `"goroutine reddit started"`, `"reddit done"`). Tu vas **voir** la concurrence se dérouler.
- Avant chaque étape, lis 5 minutes la doc Go officielle (effective_go, ou tour of Go) sur le concept que tu vas utiliser.
- À la fin de chaque étape, demande-toi : _"qu'est-ce que ça aurait donné en TS avec Promise/AbortController ?"_ — ça te fera des ponts mentaux solides.

### Quand venir me voir

Tu codes, et tu reviens me parler **uniquement quand** :

- Tu bloques sur un concept et tu veux qu'on en discute (sans code de ma part)
- Tu veux que je review ton code une fois une étape finie
- Tu veux qu'on debug ensemble un truc bizarre
- Tu as une question théorique
