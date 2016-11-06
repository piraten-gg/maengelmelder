var app = function() {
  var self = {};

  var defaultCoords = [49.9008, 8.3500];

  // show/hide fields for adding a mail address to subscribe to a new issue
  function updateSubcribeForm() {
    if($('#checkbox_subscribe').prop("checked")) {
      $('.subscribebox').show();
    } else {
      $('.subscribebox').hide();
    }
  }
  updateSubcribeForm();
  $('#checkbox_subscribe').on('change', updateSubcribeForm);

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
  var markerPos = {lat: -1, lon: -1, zoom: -1};

  // set the location of a new marker
  function setLocation(lat, lon, zoom) {
    markerPos.lat = lat;
    markerPos.lon = lon;
    markerPos.zoom = zoom;

    var baseURL = "https://nominatim.openstreetmap.org/reverse?format=json";
    $.getJSON(
      baseURL+"&lat="+lat+"&lon="+lon+"&zoom="+zoom+"&addressdetails=1",
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
      marker = L.marker([lat, lon], {draggable: 'true', icon: blueIcon}).addTo(map);

      marker.on('dragend', function(e){
        var position = marker.getLatLng();
        setLocation(position.lat, position.lng, map.getZoom());
      });
    }

    marker.setLatLng(new L.LatLng(lat, lon));
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

  self.submitMarker = function() {
    var newMarker = {};
    newMarker.lat = markerPos.lat;
    newMarker.lon = markerPos.lon;
    newMarker.zoom = markerPos.zoom;
    newMarker.category = parseInt($('.new-marker-form .category').val()) || -1;
    newMarker.descrption = $('.new-marker-form .description').val() || "";
    newMarker.confidential = $('#checkbox_confidential').is(':checked');
    newMarker.user_name = $('.new-marker-form .subName').val() || "";
    newMarker.user_mail = $('.new-marker-form .subMail').val() || "";

    if (newMarker.lat == -1) {
      alert("Bitte wählen Sie zuerst einen Punkt auf der Karte aus.")
    }

    $.post('/api/v1/markers/new', JSON.stringify(newMarker), 'json')
      .done(function(data){
        $('.new-marker-form')[0].reset();
        $('.new-marker-form .location').text('Bitte den Ort auf der Karte auswählen.');
      })
      .fail(function(data){
        console.log("failed to create marker", data);
      });

    console.log(newMarker);
  }

  return self;
}();

// {"lon": 3.4, "lat": 2.2, "category": 3, "desc": "hier ist es", "confidential": false, "user_name":"-", "user_mail":"-"}
