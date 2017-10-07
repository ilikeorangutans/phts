import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from "./dashboard/dashboard.component";
import { CollectionComponent } from "./collections/collection/collection.component";
import { CollectionDashboardComponent } from "./collections/collection-dashboard/collection-dashboard.component";
import { CollectionCreateComponent } from "./collections/collection-create/collection-create.component";

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
            }
        ]
    }
    // },
    // {
    //     path: 'collections',
    //     component: CollectionsComponent,
    //     children: [
    //         {
    //             path: '',
    //             component: DashboardComponent // TODO this should be list collections?
    //         },
    //         {
    //             path: ':slug',
    //             component: DetailsComponent,
    //             children: [
    //                 {
    //                     path: 'photos/:photoID',
    //                     component: PhotoDetailsComponent
    //                 }
    //             ]
    //         },
    //     ]
    // }
    // {
    //     path: 'collections/:slug',
    //     component: DetailsComponent,
    //     children: [
    //         {
    //             path: 'photos/:photoID',
    //             component: PhotoDetailsComponent
    //         }
    //     ]
    // }
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutesModule {}
