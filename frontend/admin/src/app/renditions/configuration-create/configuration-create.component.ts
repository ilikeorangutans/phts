import { RenditionConfigurationService } from './../../services/rendition-configuration.service';
import { RenditionConfiguration } from './../../models/rendition-configuration';
import { Collection } from './../../models/collection';
import { CurrentCollectionService } from './../../collections/current-collection.service';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-configuration-create',
  templateUrl: './configuration-create.component.html',
  styleUrls: ['./configuration-create.component.css']
})
export class ConfigurationCreateComponent implements OnInit {

  collection: Collection;
  config: RenditionConfiguration = new RenditionConfiguration();

  constructor(
    private currentCollectionService: CurrentCollectionService,
    private renditionConfigurationService: RenditionConfigurationService
  ) {
    this
      .currentCollectionService
      .current$
      .subscribe(c => this.collection = c);
  }

  ngOnInit() {
  }

  onSubmit() {
    this.config.resize = true;

    this.renditionConfigurationService.save(this.collection, this.config)
      .then(c => console.log(c));

  }

}
