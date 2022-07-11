<script lang="ts">
  import { Loader } from '@googlemaps/js-api-loader'
  import { gps } from './stores'

  let map: google.maps.Map
  let marker: google.maps.Marker
  $: latitude = parseFloat($gps.split(',')[0])
  $: longitude = parseFloat($gps.split(',')[1])

  const loader = new Loader({
    apiKey: 'AIzaSyB_4iMiaVb00W0Dsqflh2iwYCPmjGKU9KA',
    version: 'weekly',
  })
  loader.load().then(() => {
    const ll = { lat: latitude, lng: longitude }
    map = new google.maps.Map(document.getElementById('map'), { zoom: 8, center: ll })
    marker = new google.maps.Marker({ map, position: ll, draggable: true })
    marker.addListener('dragend', onDragEnd)
  })

  function onDragEnd(event: google.maps.MapMouseEvent) {
    const { lat, lng } = marker.getPosition()
    $gps = `${lat()}, ${lng()}` // if the marker is dragged, drop any altitude
  }

  function moveMarker(lat: number, lng: number) {
    if (!map) return
    map.panTo({ lat, lng })
    marker.setPosition({ lat, lng })
  }
  $: moveMarker(latitude, longitude)
</script>

<div id="map" />
