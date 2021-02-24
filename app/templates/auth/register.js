//メールアドレス形式のチェックおよび登録可能かチェック
$(function () {
  $("#email").on("input", function () {
    //validation OKまで新規登録ボタンを無効化
    var button = $("#register-btn");
    button.attr("disabled", true);
    var email_error = false;
    //メールアドレスの形式チェック
    if (
      !/^[A-Za-z0-9]{1}[A-Za-z0-9_.-]*@{1}[A-Za-z0-9_.-]{1,}\.[A-Za-z0-9]{1,}$/.test(
        $(this).val()
      )
    ) {
      email_error = true;
    }
    if (email_error) {
      if (!$(this).nextAll("span.error-info").length) {
        $(this).after(
          '<span class="text-danger error-info">有効なメールアドレスではありません。</span>'
        );
      }
    } else {
      //メールアドレス形式が問題ない場合、DB内にすでにアドレスが存在しないかチェック
      $.ajax({
        url: "/check_email", // 通信先のURL
        type: "POST", // 使用するHTTPメソッド
        data: JSON.stringify({ email: $(this).val() }),
        contentType: "application/json",
        dataType: "json", // responseのデータの種類
        timespan: 1000, // 通信のタイムアウトの設定(ミリ秒)
      })
        //通信成功
        .done(function (data, textStatus, jqXHR) {
          //エラ〜メッセージを消す
          if ($("#email").nextAll("span.error-info").length) {
            $("#email").nextAll("span.error-info").remove();
          }
          if (data.valid === true) {
            button.attr("disabled", false); // ボタンを再び enableにする
          } else {
            $("#email").after(
              '<span class="text-danger error-info">このメールアドレスはすでに登録されています。</span>'
            );
          }
        })
        //通信失敗
        .fail(function (xhr, status, error) {
          // HTTPエラー時
          //通信終了後
        });
    }
  });
});

//パスワードの文字数チェック
$(function () {
  $("#password").on("input", function () {
    //validation OKまで新規登録ボタンを無効化
    var button = $("#register-btn");
    button.attr("disabled", true);
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
      button.attr("disabled", false);
      if ($(this).nextAll("span.error-info").length) {
        $(this).nextAll("span.error-info").remove();
      }
    }
  });
});

//確認用パスワードチェック
$(function () {
  $("#password-confirm").on("input", function () {
    //validation OKまで新規登録ボタンを無効化
    var button = $("#register-btn");
    button.attr("disabled", true);
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
