import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { NotFoundComponent } from './not-found/not-found.component';
import { LandingPageComponent } from './landing-page/landing-page.component';
import { ShareViewerComponent } from './share-viewer/share-viewer.component';

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    component: LandingPageComponent,
  },
  {
    path: 'share/:slug',
    component: ShareViewerComponent,
  },
  {
    path: 'share/:slug/:index',
    component: ShareViewerComponent,
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
