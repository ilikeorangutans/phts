import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { AlbumDetailsComponent } from './album-details/album-details.component';
import { AlbumsDashboardComponent } from './albums-dashboard/albums-dashboard.component';
import { PhotoStreamComponent } from './photo-stream/photo-stream.component';
import { RenditionConfigurationsComponent } from './rendition-configurations/rendition-configurations.component';
import { LandingComponent } from './landing/landing.component';
import { AppComponent } from './app/app.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { PhotoDetailsComponent } from './photo-details/photo-details.component';

const routes: Routes = [
  {
    path: '',
    component: AppComponent,
    children: [
      {
        path: '',
        pathMatch: 'full',
        component: DashboardComponent
      }
    ]
  },
  {
    path: ':slug',
    component: AppComponent,
    children: [
      {
        path: '',
        pathMatch: 'full',
        component: LandingComponent
      },
      {
        path: 'rendition-configurations',
        component: RenditionConfigurationsComponent
      },
      {
        path: 'albums',
        component: AlbumsDashboardComponent,
      },
      {
        path: 'albums/:album_slug',
        component: AlbumDetailsComponent
      },
      {
        path: 'albums/:album_slug/photos',
        component: AlbumDetailsComponent
      },
      {
        path: 'photos',
        component: PhotoStreamComponent
      },
      {
        path: 'photos/:photo_id',
        component: PhotoDetailsComponent
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CollectionRoutingModule { }
