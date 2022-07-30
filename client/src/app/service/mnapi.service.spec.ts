import { TestBed } from '@angular/core/testing';

import { MNAPIService } from './mnapi.service';

describe('MNAPIService', () => {
  let service: MNAPIService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MNAPIService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
