import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/switchMap';
import { ActivatedRoute, ParamMap, Params } from '@angular/router';
import { CollectionService } from '../../services/collection.service';
import { RenditionConfigurationService } from "../../services/rendition-configuration.service";
import { PhotoService } from '../../services/photo.service';
import { Collection } from '../../models/collection';
import { Rendition } from '../../models/rendition';
import { Photo } from '../../models/photo';
import { RenditionConfiguration } from '../../models/rendition-configuration';
import { PathService } from '../../services/path.service';

@Component({
  selector: 'app-photo-details',
  templateUrl: './photo-details.component.html',
  styleUrls: ['./photo-details.component.css']
})
export class PhotoDetailsComponent implements OnInit {

  loading: boolean = true;
  photo: Photo;
  photoID: number;
  collection: Collection;
  renditionConfigurations: RenditionConfiguration[];
  preview: Rendition;

  constructor(
    private route: ActivatedRoute,
    private collectionService: CollectionService,
    private photoService: PhotoService,
    private renditionConfigurationService: RenditionConfigurationService,
    private pathService: PathService
  ) { }

  ngOnInit() {
    this.route
      .paramMap
      .switchMap((params: ParamMap) => {
        this.loading = true;
        this.photoID = +params.get('photoID');
        return this.collectionService.bySlug(params.get('slug'));
      })
      .subscribe(collection => this.loadPhoto(collection))
  }

  loadPhoto(collection: Collection) {
    this.collection = collection;
    console.log(`Loading photo ${this.photoID} of collection ${collection.id}`)

    this.photoService.byID(collection, this.photoID, [])
      .then((photo) => {
        console.log("loaded photo");
        this.photo = photo;

        this.preview = photo.renditions.find(r => r.renditionConfigurationID === 3)
        this.loading = false;
      })
      .catch((e) => {
        console.log(e);
      });
  }
}
