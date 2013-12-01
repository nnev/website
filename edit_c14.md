---
layout: default
title: c¼h editieren
---
<form method="POST" action="edit_c14.html">
	<input type="hidden" name="id" value="<<.Id>>"></input>

	<label for="speaker">Vortragender</label>
	<input type="text" placeholder="Nick" name="speaker" value="<<.Speaker>>"></input><br>

	<label for="topic">Thema</label>
	<input type="text" placeholder="Thema" name="topic" value="<<.Topic>>"></input><br>

	<label for="date">Datum</label>
	<input type="date" placeholder="Datum (YYYY-MM-DD)" name="date" value="<<if .HasDate >><< .Date >><<end>>"></input><br>

	<textarea name="abstract" Placeholder="Zusammenfassung" rows="5" cols="40"><<.Abstract>></textarea><br>

	<input type="submit" value="Submit"></input>

</form>
