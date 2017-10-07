import { TestBed, inject } from '@angular/core/testing';

import { CurrentCollectionService } from './current-collection.service';

describe('CurrentCollectionService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [CurrentCollectionService]
    });
  });

  it('should be created', inject([CurrentCollectionService], (service: CurrentCollectionService) => {
    expect(service).toBeTruthy();
  }));
});
