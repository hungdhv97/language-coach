.
├── index.html
├── package.json
├── tsconfig.json
├── tsconfig.app.json
├── tsconfig.node.json
├── vite.config.ts
├── public/
│   ├── favicon.ico
│   ├── robots.txt
│   ├── locales/
│   │   ├── en/
│   │   │   └── common.json
│   │   └── vi/
│   │       └── common.json
│   └── fonts/
│       └── inter/
│           ├── inter-regular.woff2
│           └── inter-medium.woff2
└── src/
    ├── main.tsx
    ├── vite-env.d.ts
    ├── app/
    │   ├── App.tsx
    │   ├── layout/
    │   │   ├── RootLayout.tsx
    │   │   ├── AuthLayout.tsx
    │   │   └── DashboardLayout.tsx
    │   ├── router/
    │   │   ├── index.tsx
    │   │   ├── routes.tsx
    │   │   ├── route-ids.ts
    │   │   └── guards/
    │   │       ├── AuthGuard.tsx
    │   │       └── RoleGuard.tsx
    │   ├── providers/
    │   │   ├── AppProviders.tsx
    │   │   ├── QueryClientProvider.tsx
    │   │   ├── ThemeProvider.tsx
    │   │   ├── I18nProvider.tsx
    │   │   └── StoreProvider.tsx
    │   └── store/
    │       ├── useAppStore.ts
    │       └── useUiStore.ts
    │
    ├── shared/
    │   ├── api/
    │   │   ├── http-client.ts
    │   │   ├── config.ts
    │   │   ├── interceptors/
    │   │   │   ├── auth-interceptor.ts
    │   │   │   └── error-interceptor.ts
    │   │   └── query/
    │   │       ├── query-client.ts
    │   │       ├── query-keys.ts
    │   │       └── types.ts
    │   ├── config/
    │   │   ├── env.ts
    │   │   ├── app-config.ts
    │   │   └── permissions.ts
    │   ├── constants/
    │   │   ├── date.ts
    │   │   ├── pagination.ts
    │   │   └── user.ts
    │   ├── lib/
    │   │   ├── utils.ts
    │   │   ├── format/
    │   │   │   ├── format-date.ts
    │   │   │   ├── format-number.ts
    │   │   │   └── format-phone.ts
    │   │   ├── i18n/
    │   │   │   ├── i18n.ts
    │   │   │   └── resources.ts
    │   │   └── validators/
    │   │       ├── zod-utils.ts
    │   │       └── rules.ts
    │   ├── hooks/
    │   │   ├── useDebounce.ts
    │   │   ├── useToggle.ts
    │   │   ├── usePagination.ts
    │   │   ├── useDisclosure.ts
    │   │   ├── useSearchParamsState.ts
    │   │   └── useBreakpoint.ts
    │   ├── layout/
    │   │   ├── AppHeader.tsx
    │   │   ├── AppSidebar.tsx
    │   │   ├── AppFooter.tsx
    │   │   ├── AppBreadcrumbs.tsx
    │   │   └── AppShell.tsx
    │   ├── store/
    │   │   ├── useThemeStore.ts
    │   │   ├── useSidebarStore.ts
    │   │   └── persist-config.ts
    │   ├── types/
    │   │   ├── common.ts
    │   │   ├── option.ts
    │   │   └── api.ts
    │   ├── ui/
    │   │   ├── icons/
    │   │   │   ├── AppLogo.tsx
    │   │   │   └── index.ts
    │   │   ├── feedback/
    │   │   │   ├── PageLoader.tsx
    │   │   │   ├── FullscreenSpinner.tsx
    │   │   │   ├── EmptyState.tsx
    │   │   │   └── ErrorState.tsx
    │   │   ├── layout/
    │   │   │   ├── PageHeader.tsx
    │   │   │   ├── PageActions.tsx
    │   │   │   └── PageContent.tsx
    │   │   ├── data-display/
    │   │   │   ├── DataTable.tsx
    │   │   │   ├── DataTableToolbar.tsx
    │   │   │   ├── DataTablePagination.tsx
    │   │   │   └── StatCard.tsx
    │   │   ├── forms/
    │   │   │   ├── FormSection.tsx
    │   │   │   └── FormFooter.tsx
    │   │   └── shadcn/
    │   │       └── ui/
    │   │           ├── accordion.tsx
    │   │           ├── alert-dialog.tsx
    │   │           ├── alert.tsx
    │   │           ├── avatar.tsx
    │   │           ├── badge.tsx
    │   │           ├── button.tsx
    │   │           ├── card.tsx
    │   │           ├── checkbox.tsx
    │   │           ├── dialog.tsx
    │   │           ├── dropdown-menu.tsx
    │   │           ├── form.tsx
    │   │           ├── input.tsx
    │   │           ├── label.tsx
    │   │           ├── popover.tsx
    │   │           ├── select.tsx
    │   │           ├── separator.tsx
    │   │           ├── sheet.tsx
    │   │           ├── skeleton.tsx
    │   │           ├── switch.tsx
    │   │           ├── table.tsx
    │   │           ├── tabs.tsx
    │   │           ├── textarea.tsx
    │   │           └── toast.tsx
    │   └── utils/
    │       ├── query-helpers.ts
    │       ├── table-helpers.ts
    │       └── form-helpers.ts
    │
    ├── forms/
    │   ├── resolvers/
    │   │   ├── zod-resolver.ts
    │   │   └── index.ts
    │   └── user/
    │       ├── useUserCreateForm.ts
    │       ├── useUserEditForm.ts
    │       ├── useUserPreferencesForm.ts
    │       ├── useChangePasswordForm.ts
    │       └── useUserInvitationForm.ts
    │
    ├── entities/
    │   └── user/
    │       ├── model/
    │       │   ├── user.types.ts
    │       │   ├── user.enums.ts
    │       │   ├── user.constants.ts
    │       │   ├── user.mappers.ts
    │       │   └── user.factories.ts
    │       ├── api/
    │       │   ├── user.api.ts
    │       │   ├── user.endpoints.ts
    │       │   ├── user.keys.ts
    │       │   ├── user.queries.ts
    │       │   └── user.mutations.ts
    │       ├── store/
    │       │   ├── useCurrentUserStore.ts
    │       │   ├── useUserFilterStore.ts
    │       │   └── useUserSessionStore.ts
    │       ├── hooks/
    │       │   ├── useCurrentUser.ts
    │       │   ├── useUserPermissions.ts
    │       │   ├── useUserSearch.ts
    │       │   └── useIsUserOnline.ts
    │       ├── ui/
    │       │   ├── UserAvatar.tsx
    │       │   ├── UserBadge.tsx
    │       │   ├── UserRoleTag.tsx
    │       │   ├── UserStatusBadge.tsx
    │       │   └── UserInlineCard.tsx
    │       └── index.ts
    │
    ├── features/
    │   └── user/
    │       ├── list/
    │       │   ├── components/
    │       │   │   ├── UserListTable.tsx
    │       │   │   ├── UserListToolbar.tsx
    │       │   │   ├── UserFilters.tsx
    │       │   │   ├── UserBulkActions.tsx
    │       │   │   └── UserListPagination.tsx
    │       │   ├── hooks/
    │       │   │   ├── useUserList.ts
    │       │   │   └── useUserListFilters.ts
    │       │   ├── store/
    │       │   │   └── useUserListStore.ts
    │       │   ├── api/
    │       │   │   └── user-list.queries.ts
    │       │   └── index.ts
    │       ├── detail/
    │       │   ├── components/
    │       │   │   ├── UserSummaryCard.tsx
    │       │   │   ├── UserProfileInfo.tsx
    │       │   │   ├── UserSecurityCard.tsx
    │       │   │   ├── UserActivityTimeline.tsx
    │       │   │   └── UserRelations.tsx
    │       │   ├── hooks/
    │       │   │   └── useUserDetail.ts
    │       │   └── index.ts
    │       ├── create/
    │       │   ├── components/
    │       │   │   ├── UserCreateForm.tsx
    │       │   │   └── UserCreateSteps.tsx
    │       │   ├── validation/
    │       │   │   └── user-create.schema.ts
    │       │   ├── hooks/
    │       │   │   └── useCreateUser.ts
    │       │   └── index.ts
    │       ├── edit/
    │       │   ├── components/
    │       │   │   ├── UserEditForm.tsx
    │       │   │   ├── UserEditTabs.tsx
    │       │   │   └── UserEditToolbar.tsx
    │       │   ├── validation/
    │       │   │   └── user-edit.schema.ts
    │       │   ├── hooks/
    │       │   │   └── useUpdateUser.ts
    │       │   └── index.ts
    │       ├── delete/
    │       │   ├── components/
    │       │   │   └── UserDeleteDialog.tsx
    │       │   ├── hooks/
    │       │   │   └── useDeleteUser.ts
    │       │   └── index.ts
    │       ├── impersonate/
    │       │   ├── components/
    │       │   │   ├── UserImpersonateBanner.tsx
    │       │   │   └── StopImpersonateButton.tsx
    │       │   ├── hooks/
    │       │   │   └── useImpersonateUser.ts
    │       │   └── index.ts
    │       ├── roles-permissions/
    │       │   ├── components/
    │       │   │   ├── UserRolesForm.tsx
    │       │   │   ├── UserPermissionsMatrix.tsx
    │       │   │   └── UserRoleAssignDialog.tsx
    │       │   ├── validation/
    │       │   │   └── user-roles.schema.ts
    │       │   ├── hooks/
    │       │   │   └── useUserRoles.ts
    │       │   └── index.ts
    │       ├── sessions/
    │       │   ├── components/
    │       │   │   ├── UserSessionsTable.tsx
    │       │   │   └── TerminateSessionDialog.tsx
    │       │   ├── hooks/
    │       │   │   └── useUserSessions.ts
    │       │   └── index.ts
    │       ├── security/
    │       │   ├── components/
    │       │   │   ├── ChangePasswordForm.tsx
    │       │   │   ├── TwoFactorSetupForm.tsx
    │       │   │   └── SecurityEventLog.tsx
    │       │   ├── validation/
    │       │   │   ├── change-password.schema.ts
    │       │   │   └── two-factor.schema.ts
    │       │   ├── hooks/
    │       │   │   └── useUserSecurity.ts
    │       │   └── index.ts
    │       ├── preferences/
    │       │   ├── components/
    │       │   │   ├── UserPreferencesForm.tsx
    │       │   │   ├── UserNotificationSettings.tsx
    │       │   │   └── UserLocalizationSettings.tsx
    │       │   ├── validation/
    │       │   │   └── user-preferences.schema.ts
    │       │   ├── hooks/
    │       │   │   └── useUserPreferences.ts
    │       │   └── index.ts
    │       ├── invitations/
    │       │   ├── components/
    │       │   │   ├── UserInvitationForm.tsx
    │       │   │   ├── UserInvitationList.tsx
    │       │   │   └── ResendInvitationDialog.tsx
    │       │   ├── validation/
    │       │   │   └── user-invitation.schema.ts
    │       │   ├── hooks/
    │       │   │   └── useUserInvitations.ts
    │       │   └── index.ts
    │       ├── import-export/
    │       │   ├── components/
    │       │   │   ├── UserImportDialog.tsx
    │       │   │   ├── UserImportMappingForm.tsx
    │       │   │   ├── UserImportPreviewTable.tsx
    │       │   │   └── UserExportDialog.tsx
    │       │   ├── hooks/
    │       │   │   ├── useImportUsers.ts
    │       │   │   └── useExportUsers.ts
    │       │   └── index.ts
    │       ├── audit-log/
    │       │   ├── components/
    │       │   │   └── UserAuditLogTable.tsx
    │       │   ├── hooks/
    │       │   │   └── useUserAuditLog.ts
    │       │   └── index.ts
    │       ├── bulk-edit/
    │       │   ├── components/
    │       │   │   ├── UserBulkEditDrawer.tsx
    │       │   │   └── UserBulkEditForm.tsx
    │       │   ├── hooks/
    │       │   │   └── useUserBulkEdit.ts
    │       │   └── index.ts
    │       ├── reporting/
    │       │   ├── components/
    │       │   │   ├── UserStatsCards.tsx
    │       │   │   └── UserReportFilters.tsx
    │       │   ├── hooks/
    │       │   │   └── useUserReports.ts
    │       │   └── index.ts
    │       ├── ui/
    │       │   ├── UserPageHeader.tsx
    │       │   ├── UserToolbar.tsx
    │       │   └── UserTabs.tsx
    │       └── index.ts
    │
    ├── pages/
    │   ├── auth/
    │   │   ├── LoginPage.tsx
    │   │   ├── RegisterPage.tsx
    │   │   ├── ForgotPasswordPage.tsx
    │   │   ├── ResetPasswordPage.tsx
    │   │   └── VerifyEmailPage.tsx
    │   ├── dashboard/
    │   │   └── DashboardHomePage.tsx
    │   └── user/
    │       ├── UserListPage.tsx
    │       ├── UserCreatePage.tsx
    │       ├── UserDetailPage.tsx
    │       ├── UserEditPage.tsx
    │       ├── UserSecurityPage.tsx
    │       ├── UserPreferencesPage.tsx
    │       ├── UserInvitationsPage.tsx
    │       └── UserAuditLogPage.tsx
    │
    └── styles/
        ├── globals.css
        ├── tailwind.css
        └── shadcn.css
