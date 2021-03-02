//有効値が入力されるまでボタンを無効化
$(function () {
  var button = $("#update-btn");
  button.attr("disabled", true);
});

//パスワードの文字数チェック
$(function () {
  $("#current-password").on("input", function () {
    //validation OKまで新規登録ボタンを無効化
    var button = $("#update-btn");
    button.attr("disabled", true);
    var password_error = false;
    if ($(this).val().length < 8) {
      password_error = true;
    }
    if (password_error) {
      if (!$(this).nextAll("span.error-info").length) {
        $(this).after(
          '<span class="text-danger error-info">パスワードが８文字以下です。</span>'
        );
      }
    } else {
      if ($(this).nextAll("span.error-info").length) {
        $(this).nextAll("span.error-info").remove();
      }
    }
  });
});

//パスワードの文字数チェック
$(function () {
  $("#password").on("input", function () {
    var password_error = false;
    if ($(this).val().length < 8) {
      password_error = true;
    }
    if (password_error) {
      if (!$(this).nextAll("span.error-info").length) {
        $(this).after(
          '<span class="text-danger error-info">パスワードは８文字以上入力してください。</span>'
        );
      }
    } else {
      if ($(this).nextAll("span.error-info").length) {
        $(this).nextAll("span.error-info").remove();
      }
    }
  });
});

//確認用パスワードチェック
$(function () {
  $("#password-confirm").on("input", function () {
    var button = $("#update-btn");
    var confirm_pass_error = false;
    if ($(this).val() !== $("#password").val()) {
      confirm_pass_error = true;
    }
    if (confirm_pass_error) {
      if (!$(this).nextAll("span.error-info").length) {
        $(this).after(
          '<span class="text-danger error-info">パスワードが一致しません。</span>'
        );
      }
    } else {
      button.attr("disabled", false);
      if ($(this).nextAll("span.error-info").length) {
        $(this).nextAll("span.error-info").remove();
      }
    }
  });
});
