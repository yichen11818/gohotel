/**
 * @see https://umijs.org/docs/max/access#access
 * */
export default function access(
  initialState: { currentUser?: API.CurrentUser } | undefined,
) {
  const { currentUser } = initialState ?? {};
  const currentUserRole =
    (currentUser as (API.CurrentUser & { role?: string; access?: string }) | undefined)?.role ||
    (currentUser as (API.CurrentUser & { role?: string; access?: string }) | undefined)?.access;
  return {
    canAdmin: currentUserRole === 'admin',
  };
}
