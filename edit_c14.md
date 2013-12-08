---
layout: default
title: c¼h editieren
---

<h1>c¼h bearbeiten/hinzufügen</h1>

<form method="POST" action="edit_c14.html">
	<input type="hidden" name="id" value="<< if ge .Id 0 >><<.Id>><< end >>" />

	<label for="speaker">Vortragender</label>
	<input type="text" placeholder="(Nick)name" id="speaker"  name="speaker" value="<<.Speaker>>" /><br>

	<label for="topic">Thema</label>
	<input type="text" placeholder="Thema" id="topic" name="topic" value="<<.Topic>>" /><br>

	<label for="date">Datum</label>
	<input type="date" placeholder="Datum (YYYY-MM-DD)" id="date" name="date" value="<<if .HasDate >><< .Date >><<end>>" /><br>

	<label for="abstract">Zusammenfassung</label>
	<textarea id="abstract" name="abstract" placeholder="Zusammenfassung" rows="10" cols="60"><<.Abstract>></textarea><br>

	<input type="submit" value="c¼h speichern"></input>

</form>
