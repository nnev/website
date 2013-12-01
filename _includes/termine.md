{% for termin in page.termine %}
<p>
	Datum: {{ termin.date }}<br>
	{% if termin.stammtisch %}
	Stammtisch<br>
	Location: {{ termin.location }}
	{% else %}
	Treff<br>
	c¼h: {{ termin.topic }}
	{% endif %}
</p>
{% endfor %}
