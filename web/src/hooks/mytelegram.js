const tg = window.Telegram.WebApp;

export function mytelegram() {
  const onClose = () => {
    tg.close();
  };

  return {
    tg,
    onClose,
    user: tg.initDataUnsafe?.user,
  };
}
