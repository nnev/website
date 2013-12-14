---
layout: default
title: Anfahrt zum Chaos-Treff Heidelberg
---

# Anfahrt


<a href="https://maps.google.de/maps?q=49.417513,8.668394&num=1&t=m&z=18" class="qrcode">
	<img src="/img/map_google.png"/><br/>
	Google Maps
</a>

<a href="http://www.openstreetmap.org/?mlat=49.41751&mlon=8.66833#map=17/49.41751/8.66833" class="qrcode"  style="margin-left: 2rem">
	<img src="/img/map_osm.png"/><br/>
	OSM.org
</a>

<address>
4. Stock, Raum 432
Im Neuenheimer Feld 368
69120 Heidelberg
</address>

Die nächste Haltestelle ist die „[Kopfklinik](http://fahrplanauskunft.vrn.de/vrn/XSLT_TRIP_REQUEST2?language=de&sessionID=0&name_destination=Neuenheim,%20Kopfklinik&type_destination=stop)“ (Bus 31, 32). Et&shy;was wei&shy;ter ent&shy;fernt die Hal&shy;te&shy;stel&shy;le „[Bun&shy;sen&shy;gym&shy;na&shy;sium](http://fahrplanauskunft.vrn.de/vrn/XSLT_TRIP_REQUEST2?language=de&sessionID=0&name_destination=Neuenheim,%20Bunsengymnasium&type_destination=stop)“ (Tram 21, 24). Par&shy;ken ist di&shy;rekt vor dem Gebäude kostenfrei möglich, sofern Platz ist. Ansonsten gibt es genug Plätze auf Otto-Meyerhof-Zentrum Parkplatz (kostenpflichtig).

<div id="map"></div>
<script>
var map = L.map('map').setView([49.41775, 8.67040], 16);
{{site.map_js}}
L.marker([{{site.treff_lat}}, {{site.treff_lon}}]).addTo(map).bindPopup("<b>Chaos-Treff</b>", { "closeButton": false }).openPopup();
L.marker([49.41910, 8.66709]).addTo(map).bindPopup("Bushaltestelle: Kopfklinik");
L.marker([49.41694, 8.67633]).addTo(map).bindPopup("Tramhaltestelle: Bunsengymnasium");
</script>
