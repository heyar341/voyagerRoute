//入力されるまでボタンを無効化
$(function () {
  var button = $("#update-btn");
  button.attr("disabled", true);
});

//ユーザー名の文字数チェック
$(function () {
  $("#username").on("input", function () {
    var button = $("#update-btn");
    if ($(this).val().length > 0) {
      button.attr("disabled", false);
    } else {
      button.attr("disabled", true);
    }
  });
});
