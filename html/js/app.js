var app = function() {
  var defaultCoords = [49.9008, 8.3500];

  // show/hide fields for adding a mail address to subscribe to a new issue
  $('#checkbox_subscribe').on('change', function(){
    if($('#checkbox_subscribe').prop("checked")) {
      $('.subscribebox').show();
    } else {
      $('.subscribebox').hide();
    }
  });

  // create the map
  var map = L.map('map', {
    zoomControl: false
  }).setView(defaultCoords, 11);

  L.control.zoom({
    position: 'topright'
  }).addTo(map);

  var Hydda_Full = L.tileLayer(
    'https://{s}.tile.openstreetmap.se/hydda/full/{z}/{x}/{y}.png', {
    attribution: 'Tiles courtesy of <a href="http://openstreetmap.se/" target="_blank">OpenStreetMap Sweden</a> &mdash; Map data &copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
  }).addTo(map);

  L.geoJson(boundaries, {style: function(feature) {
    return {
      fillColor: 'rgba(0, 0, 0, 0)',
      weight: 3,
      opacity: 1,
      color: 'orange',
      dashArray: 0,
      fillOpacity: 1
    };
  }}).addTo(map);

  var marker = null;

  // set the location of a new marker
  function setLocation(lat, lng, zoom) {
    var baseURL = "https://nominatim.openstreetmap.org/reverse?format=json";
    $.getJSON(
      baseURL+"&lat="+lat+"&lon="+lng+"&zoom="+zoom+"&addressdetails=1",
      function( data ) {
        console.log(data);
        if (data.address.county != "Kreis Groß-Gerau") {
          alert("Der gewählte Ort liegt außerhalb des Kreises Groß-Gerau!");
          return;
        }

        $(".location").text(data.display_name);
      }
    );

    if (marker === null) {
      marker = L.marker([lat, lng], {draggable: 'true', icon: blueIcon}).addTo(map);

      marker.on('dragend', function(e){
        var position = marker.getLatLng();
        setLocation(position.lat, position.lng, map.getZoom());
      });
    }
    marker.setLatLng(new L.LatLng(lat, lng));
  }

  map.on('click', function(e) {
    setLocation(e.latlng.lat, e.latlng.lng, map.getZoom());
  });

  var iconSize = [26, 34];
  var iconAnchor = [15, 34];

  var greenIcon = L.icon({
    iconUrl: 'img/marker_green.png',
    shadowUrl: 'img/marker_shadow.png',
    iconSize:   iconSize,
    iconAnchor: iconAnchor,
  });
  var orangeIcon = L.icon({
    iconUrl: 'img/marker_orange.png',
    shadowUrl: 'img/marker_shadow.png',
    iconSize:   iconSize,
    iconAnchor: iconAnchor,
  });
  var yellowIcon = L.icon({
    iconUrl: 'img/marker_yellow.png',
    shadowUrl: 'img/marker_shadow.png',
    iconSize:   iconSize,
    iconAnchor: iconAnchor,
  });
  var blueIcon = L.icon({
    iconUrl: 'img/marker_blue.png',
    shadowUrl: 'img/marker_shadow.png',
    iconSize:   iconSize,
    iconAnchor: iconAnchor,
  });

  var icons = [greenIcon, yellowIcon, orangeIcon];

  if ("geolocation" in navigator) {
    /* geolocation is available */
    navigator.geolocation.getCurrentPosition(function(position) {
      var lat = position.coords.latitude;
      var lng = position.coords.longitude;

      setLocation(lat, lng, 15);
      map.setView(new L.LatLng(lat, lng-0.005), 15, {animation: true});
    });
  }


  $.getJSON("/api/v1/markers/show", function( data ) {
    for (var i in data) {
      var d = data[i];
      console.log(d, [d.lat, d.lon]);
      L.marker([d.lat, d.lon], {
        icon: icons[Math.floor(Math.random()*3)]
      }).addTo(map);
    }
  });
}();

// {"lon": 3.4, "lat": 2.2, "category": 3, "desc": "hier ist es", "confidential": false, "user_name":"-", "user_mail":"-"}
