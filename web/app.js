document.getElementById('strava-auth').addEventListener('click', function() {
    // Redirect to the OAuth endpoint which should handle the OAuth redirect
    window.location.href = "/auth";
});