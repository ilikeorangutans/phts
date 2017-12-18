import { ShareSitesCreateComponent } from './../share-sites/share-sites-create/share-sites-create.component';
import { ShareSitesDashboardComponent } from './../share-sites/share-sites-dashboard/share-sites-dashboard.component';
import { PhotoListComponent } from './../photos/photo-list/photo-list.component';
import { ConfigurationCreateComponent } from './../renditions/configuration-create/configuration-create.component';
import { ConfigurationListComponent } from './../renditions/configuration-list/configuration-list.component';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { CollectionBrowserComponent } from './../collections/collection-browser/collection-browser.component';
import { DashboardComponent } from '../dashboard/dashboard.component';
import { CollectionComponent } from '../collections/collection/collection.component';
import { CollectionDashboardComponent } from '../collections/collection-dashboard/collection-dashboard.component';
import { CollectionCreateComponent } from '../collections/collection-create/collection-create.component';
import { PhotoDetailsComponent } from '../photos/photo-details/photo-details.component';

const routes: Routes = [
    {
        path: '',
        redirectTo: '/dashboard',
        pathMatch: 'full'
    },
    {
        path: 'dashboard',
        component: DashboardComponent
    },
    {
        path: 'collections',
        component: CollectionBrowserComponent
    },
    {
        path: 'collections/new',
        component: CollectionCreateComponent
    },
    {
        path: 'collections/:slug',
        component: CollectionComponent,
        children: [
            {
                path: '',
                component: CollectionDashboardComponent
            },
            {
                path: 'photos',
                component: PhotoListComponent
            },
            {
                path: 'photos/:photoID',
                component: PhotoDetailsComponent
            },
            {
                path: 'rendition-configurations',
                component: ConfigurationListComponent
            },
            {
                path: 'rendition-configurations/new',
                component: ConfigurationCreateComponent
            }
        ]
    },
    {
        path: 'share-sites',
        component: ShareSitesDashboardComponent
    },
    {
        path: 'share-sites/new',
        component: ShareSitesCreateComponent
    }
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutesModule {}
