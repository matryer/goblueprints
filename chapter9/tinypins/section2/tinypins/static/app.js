var sendUpdateInterval = 30/* seconds */ * 1000
var receiveUpdateInterval = 10/* seconds */ * 1000
var map, mapEl
var meMarker, otherMarkers

function initMap() {
	if (location.protocol == 'http:' && location.host != 'localhost:8080') {
		location.href = 'https://'+location.href.split('://')[1]
	}
	mapEl = document.getElementById('map')
	map = new google.maps.Map(mapEl, {
		center: {lat: 0, lng: -51},
		zoom: 3
	})
	map.setOptions({draggable: false, zoomControl: false, scrollwheel: false, disableDoubleClickZoom: true});
	if (!navigator.geolocation) {
		setError('Sorry, your browser doesn\'t support locations')
		return
	}
	otherMarkers = {}
	updateMyPosition()
	loadOtherPins()
}

function updateMyPosition() {
	console.info('getting latest position...')
	navigator.geolocation.getCurrentPosition(onPositionUpdate, locationFailed, {
		enableHighAccuracy: true
	})
}

var lastLat, lastLng, lastAccuracy
function onPositionUpdate(position) {
	// send this to the server

	if (lastLat == position.coords.latitude
		&& lastLng == position.coords.longitude &&
			lastAccuracy == position.coords.accuracy) {
		console.info('skipping update - nothing changed')
		return
	}

	lastLat = position.coords.latitude
	lastLng = position.coords.longitude
	lastAccuracy = position.coords.accuracy
	data = {
		location: {
			lat: position.coords.latitude,
			lng: position.coords.longitude
		},
		accuracy: position.coords.accuracy
	}
	console.info('sending update', data)
	$.ajax({
		type: 'POST',
		url: location.pathname+'/update',
		dataType: 'json',
		data: JSON.stringify(data),
		success: function(response){
			console.info('(update sent to server)')
			// schedule next update in sendUpdateInterval seconds
			console.info('waiting', (sendUpdateInterval/1000)+' seconds to send update again...')
			setTimeout(updateMyPosition, sendUpdateInterval)

			if (!meMarker) {
				meMarker = new google.maps.Marker({
					map: map,
					icon: {
						url: response.image_url,
						scaledSize: new google.maps.Size(32, 32)
					},
					title: response.me ? 'You' : response.name,
					position: {
						lat: response.location.Lat,
						lng: response.location.Lng
					}
				})
				updateMapZoomAndBounds()
			}

		},
		error: function(){
			console.warn('failed to update server:', arguments)
		}
	})
}

function locationFailed(error) {
	console.warn('getCurrentPosition:', arguments)
	setError(error)
}

function loadOtherPins() {
	$.ajax({
		type: 'GET',
		url: location.pathname+'/update',
		dataType: 'json',
		success: function(pins) {
			processPins(pins)
			console.info('waiting', (receiveUpdateInterval/1000)+' seconds to get updates...')
			setTimeout(loadOtherPins, receiveUpdateInterval)
		},
		error: function(){
			console.warn('failed to get pins from server:', arguments)
		}
	})
}

function processPins(pins) {
	console.info('pins:', pins)
	for (var p in pins) {
		if (!pins.hasOwnProperty(p)) continue
		var pin = pins[p]
		var marker
		if (pin.me) {
			marker = meMarker
		} else {
			marker = otherMarkers[pin.user_id]
		}
		if (!marker) {
			marker = new google.maps.Marker({
				map: map,
				icon: {
					url: pin.image_url,
					scaledSize: new google.maps.Size(32, 32)
				},
				title: pin.me ? 'You' : pin.name,
				position: {
					lat: pin.location.Lat,
					lng: pin.location.Lng
				}
			})
			if (pin.me) {
				meMarker = marker
			} else {
				otherMarkers[pin.user_id] = marker
			}
		}

		var opacity = 1
		if (pin.seconds_ago > 1*60) {
			opacity = 0.8
		}
		if (pin.seconds_ago > 2*60) {
			opacity = 0.6
		}
		if (pin.seconds_ago > 5*60) {
			opacity = 0.5
		}
		if (pin.seconds_ago > 10*60) {
			opacity = 0.25
		}
		marker.setOptions({'opacity': opacity})
		updateMapZoomAndBounds()
	}
}

function setError(error) {
	mapEl.innerHTML = error
}

function updateMapZoomAndBounds() {
	// now all markers are updated - make sure the map fits
	var markers = [meMarker]
	for (var k in otherMarkers) {
		markers.push(otherMarkers[k])
	}
	fitToMarkers(markers)
}

function fitToMarkers(markers) {

    var bounds = new google.maps.LatLngBounds();

    // Create bounds from markers
    var didSet = false
    for(var i in markers) {
    	var marker = markers[i]
    	if (!marker) continue
        var latlng = markers[i].getPosition();
        bounds.extend(latlng)
        didSet = true
    }

    if (!didSet) {
    	return
    }

   var extendPoint1 = new google.maps.LatLng(bounds.getNorthEast().lat() + 0.003, bounds.getNorthEast().lng() + 0.003);
   var extendPoint2 = new google.maps.LatLng(bounds.getNorthEast().lat() - 0.003, bounds.getNorthEast().lng() - 0.003);
   bounds.extend(extendPoint1);
   bounds.extend(extendPoint2);

    map.fitBounds(bounds);

    // Adjusting zoom here doesn't work :/

}