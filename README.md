# AOTY Playlist generator (2.0)

Generate playlists compiled from the most talked about albums over at aoty.org.

Customize your playlist through several options such as genre, track length, and album diversity. A preview of each playlist type is given. Logging in using Spotify requires your account to be whitelisted since the app is not in extended quota mode. Generated playlists are linked and can be added to your library.

### Work in progress:

2.0 involves switching the backend from node.js to Go to utilize preformance optimizations from a compiled language as well as go routines (concurrency). Caching and containerization has been implemented for better performance and ease of deployment. This transition is still in progress and will be deployed soon. Upcoming features include generating playlists based on decade time periods, and server queue implementations for better optimization.
