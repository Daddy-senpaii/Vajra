export type MockAuthResult = {ok: false; message: string};

export async function mockSignUp(): Promise<MockAuthResult> {
  return {ok: false, message: 'Account creation is a UI preview. Authentication is not connected yet.'};
}

export async function mockGoogleAuth(): Promise<MockAuthResult> {
  return {ok: false, message: 'Google authentication is not connected yet.'};
}
