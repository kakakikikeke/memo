<section class="hero is-primary is-fullheight">
  <div class="hero-head">
    <header class="nav">
      <div class="container">
        <div class="nav-left">
          <a class="nav-item" href="/">
            memo
          </a>
        </div>
        <span class="nav-toggle">
          <span></span>
          <span></span>
          <span></span>
        </span>
        <div class="nav-right nav-menu">
          <span class="nav-item">
            <a class="button is-primary is-inverted" href="https://github.com/kakakikikeke/memo/tree/ver_golang" target="_blank">
              <span class="icon">
                <i class="fa fa-github"></i>
              </span>
              <span>Show code</span>
            </a>
          </span>
        </div>
      </div>
    </header>
  </div>
  <div class="hero-body">
    <div class="container has-text-centered">
      <h1 class="title">memo</h1>
      <h2 class="subtitle">
        Free note for Everyone.
      </h2>
      <div class="field">
        <p class="control">
          <label class="checkbox">
            <input type="checkbox" id="toggle">Textbox</input>
          </label>
        </p>
      </div>
      <div class="field">
        <p class="control">
          <input class="input" type="text" placeholder="something" id="value">
	</p>
        <p class="control">
          <textarea class="textarea" placeholder="some lines" id="values" style="display:none"></textarea>
        </p>
      </div>
      <div class="field">
        <p class="control">
          <button class="button is-info" id="submit" disabled>memo</button>
        </p>
      </div>
      {{ range $memo := .memos }}
        {{ $replaced_memo := rep $memo "\n" "<br>" }}
        <div>{{ str2html $replaced_memo }}</div>
      {{ end }}
    </div>
  </div>
  <div class="hero-foot">
    <nav class="tabs is-centered">
      <div class="container">
        <a class="button is-danger" id="clear">Clear</a>
      </div>
    </nav>
  </div>
</section>