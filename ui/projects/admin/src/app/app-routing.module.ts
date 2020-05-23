import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { NotFoundComponent } from './not-found/not-found.component';
import { JoinComponent } from './join/join.component';
import { LoginComponent } from './login/login.component';
import { AdminShellComponent } from './admin-shell/admin-shell.component';
import { AuthGuard } from './services/auth.guard';
import { DashboardComponent } from './dashboard/dashboard.component';

const routes: Routes = [
  {
    path: 'join/:token',
    component: JoinComponent,
  },
  {
    path: '',
    pathMatch: 'prefix',
    component: AdminShellComponent,
    canActivate: [AuthGuard],
    canActivateChild: [AuthGuard],
    children: [
      {
        path: 'dashboard',
        component: DashboardComponent,
      },
      {
        path: 'account',
        loadChildren: () =>
          import('./account/account.module').then((m) => m.AccountModule),
      },
      {
        path: 'collection',
        loadChildren: () =>
          import('./collection/collection.module').then(
            (m) => m.CollectionModule
          ),
      },
      {
        path: 'share-site',
        loadChildren: () =>
          import('./share-site/share-site.module').then(
            (m) => m.ShareSiteModule
          ),
      },
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'dashboard',
      },
    ],
  },
  {
    path: '',
    pathMatch: 'full',
    redirectTo: '/login',
  },
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: '**',
    component: NotFoundComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
