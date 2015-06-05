---
layout: default
title: c¼h editieren
---

<h1>c¼h bearbeiten/hinzufügen</h1>

(( if ge .Id 0 ))
<p>
	Bitte speichere folgenden Link. Du brauchst ihn, wenn du in Zukunft deinen
	Eintrag editieren willst:<br />
	<a href="/edit_c14.html?id=(( .Id ))&pw=(( .Password.String ))">http://www.noname-ev.de/edit_c14.html?id=(( .Id ))&pw=(( .Password.String ))</a>
</p>
(( end ))

<form method="POST" action="edit_c14.html">
	<input type="hidden" name="id" value="(( if ge .Id 0 ))((.Id))(( end ))" />
	<input type="hidden" name="pw" value="(( .Password.String ))" />

	<label for="speaker">Vortragender</label>
	<input type="text" placeholder="(Nick)name" id="speaker"  name="speaker" value="((.Speaker))" required="required"/><br>

	<label for="topic">Thema</label>
	<input type="text" placeholder="Thema" id="topic" name="topic" value="((.Topic))" required="required" /><br>

	<label for="date">Datum</label>
	<input type="date" placeholder="Datum (YYYY-MM-DD)" id="date" name="date" value="((if .HasDate ))(( .Date ))((end))" /><br>

	<label for="abstract">Zusammenfassung</label>
	<textarea id="abstract" name="abstract" placeholder="Zusammenfassung" rows="10" cols="60">((.Abstract))</textarea><br>

	<script>
	function addField() {
		document.getElementById("links").innerHTML +=
		'<label for="kind">Art</label>' +
		'<input type="text" placeholder="Art" id="kind" name="kind" value="" /><br>' +
		'<label for="url">Url</label>' +
		'<input type="text" placeholder="http://example.com/folien.pdf" id="url" name="url" value="" /><br>';
	}
	</script>

	<a class="button" onclick="addField()" > Informationen/Links hinzufügen</a>

	<label for="links">Informationen/Links </label>

	<div id="links">
	(( range .Links ))
	<label for="kind">Art</label>
	<input type="text" placeholder="Art" id="kind" name="kind" value="((.Kind))" /><br>
	<label for="url">Url</label>
	<input type="text" placeholder="http://example.com/folien.pdf" id="url" name="url" value="((.Url))" /><br>
	(( end ))
	</div>

	<input type="submit" value="c¼h speichern" />
	((if ge .Id 0))
	<input type="submit" name="delete" value="c¼h löschen" />
	((end))
</form>


