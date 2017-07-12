<script type="text/javascript">
$(function() {
  $(document).ready(function() {
    $('#value').val("");
  });
  $('#submit').click(function() {
    $('#submit').addClass('is-loading');
    var value = $('#value').val();
    $.ajax({
      type: "POST",
      url: "/insert",
      data: {
        "msg": value
      },
      success: function(result){
        location.reload();
      }
    });
  });
  $('#clear').click(function() {
    $('#clear').addClass('is-loading');
    $.ajax({
      type: "POST",
      url: "/clear",
      success: function(result) {
        location.reload();
      }
    });
  });
  $('#value').keyup(function() {
    if($(this).val().length != 0) {
      $('#submit').attr('disabled', false);
    } else {
      $('#submit').attr('disabled', true);
    }
  });
  $('#value').keypress(function(e) {
    if (e.which == 13)  {
      $('#submit').click();
    }
  });
});
</script>