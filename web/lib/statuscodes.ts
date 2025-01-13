const oauthError =
  "Something wrong happened during the OAuth process, try again and if the issue presists please contact the developer!";

export const statusCodesMap: { [key: number]: string } = {
  0: "No user with that data found!",
  1: "OAuth login detected, Please login with google! (if you want to set a password you can do so in the settings).",
  2: "Invalid email or password!",
  3: "An account with that email already exist!",
  4: oauthError,
  5: oauthError,
  6: oauthError,
  7: "Internal server error, Please try again!",
  8: "You're already logged in! Please logout first.",
};
