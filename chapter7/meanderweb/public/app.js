(function($){
  $(function(){

    var apiHost = "http://localhost:8080";
    var map = null;
    var markers = [];

    function busy(b) {
      if (b) {
        $(".loader").slideDown();
      } else {
        window.setTimeout(function(){$(".loader").slideUp()}, 500);
      }
    }

    $(".change-adventure").click(function(){
      $("#adventure").slideUp();
      $("#meanderform").slideDown();
    });

    $(".shuffle").click(function(){
      $("#meanderform").submit();
    });

    $("#meanderform").submit(function(e){
      e.preventDefault();
      busy(true);

      var gc = new google.maps.Geocoder();
      var locationText = $("#location").val();
      gc.geocode({
        address: locationText
      }, function(results, status){

        if (status != google.maps.GeocoderStatus.OK) {
          alert("Ahh, we don't know where \"" + locationText + "\" is.");
          busy(false);
          return;
        }

        var result = results[0];
        var location = result.geometry.location;

        $.ajax({
          url: apiHost + "/recommendations",
          dataType: "json",
          data: {
            "lat": location.lat(),
            "lng": location.lng(),
            "radius": 5000,
            "cost": $("#price").val(),
            "journey": $("#journeys").val()
          },
          success: function(results){

            var title = $('#journeys option:selected').text() + " in " + locationText;
            $("#adventure .panel-title").text(title);
            $("#adventure").slideDown();
            $("#meanderform").slideUp();

            if (!map) {
              map = new google.maps.Map($("#map")[0], {
                center: location,
                zoom: 1
              });
            }

            console.info('set center')
            map.setCenter(location);
            map.setZoom(1);

            // clear markers
            for (var m in markers) {
              markers[m].setMap(null);
              markers[m] = null;
            }
            markers = [];

            // build route
            var route = $("#route").empty();
            $("#photos").empty();
            var counter = 0;
            var bounds = new google.maps.LatLngBounds();

            for (var r in results) {
              if (!results.hasOwnProperty(r)) continue;

              var item = results[r];
              if (!item) continue;
              route.append(
                $("<li/>").append(
                  $("<span>").append(
                    $("<img>").addClass("pull-right icon").attr("src", item.icon).css({width:20,height:20}),
                    " ",
                    item.name,
                    " ",
                    $("<small>").text(item.vicinity)
                  )
                )
              )

              if (item.photos) {
                if (item.photos.length > 0) {
                  console.info(item.photos)
                  $("#photos").append(
                    $("<li/>").append(
                      $("<img>").attr("src", item.photos[0].url)
                    )
                  );
                }
              }

              var thisLocation = new google.maps.LatLng(item.geometry.location.lat,item.geometry.location.lng);
              bounds.extend(thisLocation);

              window.setTimeout((function(){
                var $thisLocation = thisLocation;
                var $item = item;
                var $counter = counter+1;
              return function(){
                  var m = new google.maps.Marker({
                    position: $thisLocation,
                    map: map,
                    title: $item.name,
                    animation: google.maps.Animation.DROP,
                    icon: "https://chart.googleapis.com/chart?chst=d_map_pin_letter&chld=" + $counter + "|FF776B|000000"
                    //icon: new google.maps.MarkerImage($item.icon, null, null, null, new google.maps.Size(25, 25))
                  });
                  markers.push(m);
                }
              })(), 500*(counter++))

            }

            map.fitBounds(bounds);

          },
          complete: function(){
            busy(false);
          }
        });

      })

    });

    // load the journeys
    $.ajax({
      type: "GET",
      url: apiHost + "/journeys",
      error: function() {
        alert("The API doesn't seem to be running on :8080");
        busy(false);
      },
      dataType: "json",
      success: function(journeys) {

        for (var i in journeys) {
          var journey = journeys[i];
          $("#journeys").append(
            $("<option/>")
              .text(journey.name)
              .val(journey.journey)
          )
        }

        busy(false);

      }
    });

  });
})(jQuery);