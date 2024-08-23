document.getElementById('strava-auth').addEventListener('click', function() {
    // Redirect to the root path which should handle the OAuth redirect
    window.location.href = "/";
});

// You can also add logic to load and initialize the WebAssembly module here.