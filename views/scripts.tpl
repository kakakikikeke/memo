<script type="text/javascript">
$(function() {
  $(document).ready(function() {
    $('#value,#values').val("");
  });
  $('#submit').click(function() {
    $('#submit').addClass('is-loading');
    var textType = "#value"
    if($('#toggle').is(':checked')) {
      textType = "#values"
    }
    var value = $(textType).val();
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
  $('#value,#values').keyup(function() {
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
  $('#toggle').click(function() {
    if($('#toggle').is(':checked')) {
      $('#value').css('display', 'none');
      $('#values').removeAttr('style');
    } else {
      $('#value').removeAttr('style');
      $('#values').css('display', 'none');
    }
  });
});
</script>