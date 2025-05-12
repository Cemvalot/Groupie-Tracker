document.addEventListener("DOMContentLoaded", function () {
    // Initialize the map with a default view.
    var map = L.map('map').setView([0, 0], 2);
  
    // Add the OpenStreetMap tile layer.
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);
  
    // Data structures for markers and the connecting line.
    var markers = {}; // Stores markers by location name
    var pathCoordinates = [];
    var polyline = L.polyline(pathCoordinates, { color: 'blue' }).addTo(map);
  
    // Function to display the zoom buttons.
    function showButtons() {
      document.getElementById('zoom-out-button').style.display = 'block';
      document.getElementById('zoom-in-button').style.display = 'block';
    }
  
    // Function to zoom to specific coordinates.
    function zoomToCoordinate(lat, lon) {
      map.setView([lat, lon], 20);
    }
  
    // Select all location cards.
    var cards = document.querySelectorAll('.location-card');
    cards.forEach(function (cardElement) {
      var locationName = cardElement.querySelector('.location-name').getAttribute('data-location');
      var [city, country] = locationName.split('-');
      var query = `${city}, ${country}`;
  
      // Directly fetch coordinates using the Nominatim API.
      fetch(`https://nominatim.openstreetmap.org/search?format=json&addressdetails=1&q=${encodeURIComponent(query)}`)
        .then(response => response.json())
        .then(data => {
          if (data && data.length > 0) {
            const firstResult = data[0];
            var lat = parseFloat(firstResult.lat);
            var lon = parseFloat(firstResult.lon);
  
            // Create and add a marker to the map.
            var marker = L.marker([lat, lon]).addTo(map);
            marker.bindPopup(`<strong>${locationName}</strong>`);
            markers[locationName] = marker;
            pathCoordinates.push([lat, lon]);
            cardElement.classList.add('highlight');
  
            // Update the line connecting the markers.
            // polyline.setLatLngs(pathCoordinates);
  
            // Display the zoom buttons.
            showButtons();
          } else {
            alert("Location not found: " + locationName);
          }
        })
        .catch(err => console.error("Error fetching location:", err));
  
      // Add event listener for zoom when the card is clicked.
      cardElement.addEventListener('click', function () {
        if (markers[locationName]) {
          var markerLatLng = markers[locationName].getLatLng();
          zoomToCoordinate(markerLatLng.lat, markerLatLng.lng);
        } else {
          alert("The marker has not loaded yet for " + locationName);
        }
      });
    });
  
    // Zoom Out Button: Reset the map view.
    document.getElementById('zoom-out-button').addEventListener('click', function () {
      map.setView([0, 0], 2);
    });
  
    // Zoom In Button: Zoom to the last added point.
    document.getElementById('zoom-in-button').addEventListener('click', function () {
      if (pathCoordinates.length > 0) {
        var lastCoordinate = pathCoordinates[pathCoordinates.length - 1];
        map.setView(lastCoordinate, 20);
      } else {
        alert("No location has been added for zoom.");
      }
    });
  });