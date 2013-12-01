{% for termin in page.termine %}
<p>
	Datum: {{ termin.date }}<br>
	{% if termin.stammtisch %}
	Location: {{ termin.location }}
	{% else %}
	c¼h: {{ termin.topic }}
	{% endif %}
</p>
{% endfor %}
