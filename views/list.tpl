<form action="/insert" method="post">
  <input type="text" name="msg" value="" size="80">
  <input type="submit" value="memo"/>
</form>
<form action="/clear" method="post">
  <input type="submit" value="clear"/>
</form>
<hr/>
{{ range $memo := .memos }}
<p>{{ $memo }}</p>
{{ end }}