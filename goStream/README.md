# 🔴 Exercice final — Plateforme de stream WebRTC

## Contexte

Tu construis une plateforme où **un diffuseur** stream de l'audio/vidéo en direct vers **N auditeurs**. Le streaming utilise WebRTC en P2P direct ; ton serveur Go fait le **signaling** (la négociation initiale entre peers) et la gestion des rooms.

C'est exactement le pattern de Twitch/Clubhouse simplifié : un broadcaster, plusieurs viewers, signaling centralisé.

## Stack imposée

- Tout ce que tu as utilisé pour le chat room.
- WebRTC côté **client** uniquement (JavaScript dans le navigateur, API `RTCPeerConnection` native).
- Ton serveur Go ne touche **pas** au flux audio/vidéo. Il ne fait que du signaling via WebSocket.

## Critères de réussite

1. Lancer ton serveur Go.
2. Ouvrir une page `/broadcast` dans un navigateur → tu actives ta caméra/micro et tu deviens diffuseur d'une room.
3. Ouvrir une page `/watch/<room>` dans d'autres navigateurs/onglets → tu vois et entends le diffuseur en temps réel.
4. Quand le diffuseur quitte, les auditeurs voient un message "stream terminé".
5. Quand un auditeur quitte, le diffuseur continue normalement.

## Comprendre le rôle du signaling avant de coder

Avant d'écrire une ligne, lis sur WebRTC. Voilà le minimum à comprendre :

- Deux peers veulent se parler en P2P, mais ils ne connaissent pas leurs adresses publiques.
- Ils s'échangent des **SDP offers/answers** (description de leur connexion) et des **ICE candidates** (chemins réseau possibles) via **un serveur tiers**.
- Une fois cet échange fait via le serveur, le flux audio/vidéo passe **directement** entre eux, le serveur n'est plus dans la boucle.

Ton serveur Go est ce **serveur tiers**. Son job : transporter les messages SDP et ICE entre le diffuseur et chaque auditeur.

## Les 6 étapes

### Étape 1 — Reprend ton chat room et adapte-le

- Crée une notion de **Room** : chaque room a **un seul diffuseur** et **N auditeurs**.
- Endpoint `/ws/broadcast/<room>` → upgrade WS, marque le client comme diffuseur de la room.
- Endpoint `/ws/watch/<room>` → upgrade WS, marque le client comme auditeur de la room.
- Refuse une 2e connexion broadcast sur une room déjà occupée.
- Pas encore de WebRTC. Juste : le diffuseur peut envoyer un message texte qui est broadcast à tous les auditeurs de sa room (réutilise ton pattern hub).

### Étape 2 — Côté client : pages HTML/JS

- Une page `/broadcast/<room>` qui demande l'accès caméra/micro et affiche le flux local.
- Une page `/watch/<room>` qui prépare un élément `<video>` vide.
- Pour l'instant, ces pages ne font pas encore de WebRTC. Juste l'UI et la connexion WS.

### Étape 3 — Le signaling : transporter les SDP

C'est le cœur du projet. Côté client en JS :

- Le diffuseur, à l'arrivée d'un nouvel auditeur, crée une `RTCPeerConnection`, ajoute son flux local, génère une **offer SDP**, et l'envoie au serveur en disant "pour l'auditeur X".
- Le serveur route ce message vers l'auditeur X.
- L'auditeur reçoit l'offer, crée sa propre `RTCPeerConnection`, génère une **answer SDP**, l'envoie au serveur.
- Le serveur route l'answer vers le diffuseur.

**Côté Go** : ton hub doit savoir router des messages spécifiquement entre 2 clients identifiés, pas juste broadcast à tous. Donne un ID unique à chaque client à la connexion.

### Étape 4 — Les ICE candidates

Après l'échange SDP, chaque peer génère des **ICE candidates** en continu (chemins réseau qu'il propose). Il faut les router pareil que les SDP.

- Le diffuseur génère un candidate → il l'envoie au serveur avec "pour l'auditeur X" → serveur route → l'auditeur l'ajoute à sa PeerConnection.
- Pareil dans l'autre sens.

À la fin de cette étape, si tout marche : tu lances le diffuseur, tu ouvres un auditeur, **et tu vois la vidéo passer**. 🎉

### Étape 5 — Gestion du multi-auditeurs

C'est là que ça devient intéressant côté Go. Le diffuseur doit maintenir **N connexions WebRTC en parallèle**, une par auditeur. Donc :

- Quand un nouvel auditeur arrive, le serveur notifie le diffuseur : "nouveau peer, ID=Y".
- Le diffuseur crée une nouvelle `RTCPeerConnection` dédiée à Y et lance le signaling avec lui.
- Quand un auditeur part, le serveur notifie le diffuseur : "peer Y parti", et le diffuseur ferme sa PeerConnection pour Y.

Pattern Go à utiliser : ton hub a une **map des auditeurs par room**, et envoie des événements "auditeur connecté / déconnecté" sur le channel du diffuseur.

### Étape 6 — Annulation propre

- Quand le diffuseur quitte, **tous** les auditeurs de sa room doivent être notifiés et leur connexion WS fermée proprement.
- Utilise un `context.Context` par room : quand le diffuseur part, `cancel()` propage l'annulation à tous les auditeurs.
- Vérifie qu'il n'y a aucune fuite de goroutine avec `runtime.NumGoroutine()`.

## Bonus

- **Chat textuel en parallèle du stream** : les auditeurs peuvent envoyer des messages que tous voient. Tu réutilises littéralement ton chat room de l'exo intermédiaire en parallèle de la couche signaling.
- **Liste des rooms actives** sur la page d'accueil avec leur nombre d'auditeurs en direct.
- **Reconnexion automatique** côté client si la WS tombe.
- **Authentification simple** : token dans l'URL pour empêcher n'importe qui de devenir diffuseur d'une room.

## Pièges qui vont te ralentir (sois prêt mentalement)

- WebRTC est **chiant à débugger**. Les erreurs ne sont pas toujours explicites. Tu vas passer du temps dans la console JS du navigateur.
- L'ordre des opérations compte : ajouter le track AVANT de créer l'offer, sinon l'auditeur reçoit une connexion mais pas de vidéo.
- Sans serveur **STUN/TURN**, ça ne marche qu'en LAN. Pour tester depuis 2 réseaux différents, utilise les serveurs STUN publics de Google : `stun:stun.l.google.com:19302`. Pas besoin de TURN si tu testes en local.
- HTTPS requis pour `getUserMedia` sauf sur `localhost`. Donc teste en local d'abord.

## Quand tu auras fini

Tu auras construit, en Go pur, un système qui :

- Gère N clients connectés simultanément.
- Route des messages entre clients spécifiques.
- Maintient un état partagé (rooms) sans race condition.
- Propage des annulations proprement à travers une hiérarchie de goroutines.
- S'intègre avec une techno temps réel moderne (WebRTC).

C'est **exactement** le profil de compétences attendu pour un backend Go orienté temps réel/streaming.

---

## Conseils transverses pour les deux exercices

**Avant de coder, fais un schéma sur papier** de qui parle à qui. Pour le chat room : un cercle "Hub", des cercles "Client", flèches avec les noms des channels. Pour le projet WebRTC : pareil + le détail du signaling.

**Code par petits commits qui marchent.** Chaque étape = au moins un commit. Comme ça si tu casses tout, tu reviens en arrière sans rage.

**Mesure le nombre de goroutines** régulièrement. C'est le révélateur n°1 des bugs de concurrence en Go. Une appli qui marche mais leak des goroutines, c'est une bombe à retardement en prod.

**Ne pose de questions qu'au bout de 30 minutes de blocage minimum.** Et formule la question en termes de **concept**, pas de "mon code marche pas".

Bon courage, à dans quelques jours pour debriefer. 💪
